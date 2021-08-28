package chargepoint

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/kr/pretty"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ws"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type ChargePointHandler struct {
	chargePoint ocpp16.ChargePoint
	IsAvailable bool
	Connectors  []*Connector
	Settings    *settings.Settings
	TagReader   *hardware.TagReader
	LEDStrip    *hardware.LEDStrip
	LCD         *hardware.LCD
}

func (handler *ChargePointHandler) sendToLCD(message hardware.LCDMessage) {
	if handler.LCD == nil {
		return
	}
	handler.LCD.LCDChannel <- message
}

func (handler *ChargePointHandler) displayLEDStatus(connectorIndex int, status core.ChargePointStatus) {
	if !handler.Settings.ChargePoint.Hardware.LedIndicator.Enabled || handler.LEDStrip == nil {
		return
	}
	var color = 0x00
	switch status {
	case core.ChargePointStatusFaulted:
		color = hardware.RED
		break
	case core.ChargePointStatusCharging:
		color = hardware.BLUE
		break
	case core.ChargePointStatusReserved:
		color = hardware.YELLOW
		break
	case core.ChargePointStatusFinishing:
		color = hardware.BLUE
		break
	case core.ChargePointStatusAvailable:
		color = hardware.GREEN
		break
	case core.ChargePointStatusUnavailable:
		color = hardware.ORANGE
		break
	}
	if color != 0x00 {
		_, err := scheduler.Every(1).Milliseconds().LimitRunsTo(1).Do(handler.LEDStrip.DisplayColor, connectorIndex, uint32(color))
		if err != nil {
			log.Printf("Error display LED status: %v \n", err)
		}
	}
}

func (handler *ChargePointHandler) indicateCard(index int, color uint32) {
	if !handler.Settings.ChargePoint.Hardware.LedIndicator.Enabled || handler.LEDStrip == nil {
		return
	}
	_, err := scheduler.Every(1).Milliseconds().LimitRunsTo(1).Do(handler.LEDStrip.Blink, index, 3, color)
	if err != nil {
		log.Printf("Error indicating card: %v", err)
		return
	}
}

// FindAvailableConnector Find first Connector with the status "Available" from the handler.
func (handler *ChargePointHandler) FindAvailableConnector() *Connector {
	for _, connector := range handler.Connectors {
		if connector.IsAvailable() {
			return connector
		}
	}
	return nil
}

// FindConnectorWithId Find the Connector with the specified connectorID.
func (handler *ChargePointHandler) FindConnectorWithId(connectorID int) *Connector {
	for _, connector := range handler.Connectors {
		if connector.ConnectorId == connectorID {
			return connector
		}
	}
	return nil
}

// FindConnectorWithTagId Find the Connector that has the same tagId as the session of the connector.
func (handler *ChargePointHandler) FindConnectorWithTagId(tagId string) *Connector {
	for _, connector := range handler.Connectors {
		if connector.session.TagId == tagId {
			return connector
		}
	}
	return nil
}

// FindConnectorWithTransactionId Find the Connector that contains the transactionId in the session of the connector.
func (handler *ChargePointHandler) FindConnectorWithTransactionId(transactionId string) *Connector {
	for _, connector := range handler.Connectors {
		if connector.session.TransactionId == transactionId {
			return connector
		}
	}
	return nil
}

// FindConnectorWithReservationId Find the Connector that contains the reservationId.
func (handler *ChargePointHandler) FindConnectorWithReservationId(reservationId int) *Connector {
	for _, connector := range handler.Connectors {
		if connector.GetReservationId() == reservationId {
			return connector
		}
	}
	return nil
}

// startup After connecting to the central system, try to restore the previous state of each Connector and notify about its state.
//If the ConnectorStatus was "Preparing" or "Charging",  try to resume or start charging. If it fails, notify the central system.
func (handler *ChargePointHandler) startup() {
	var err error
	for _, connector := range handler.Connectors {
		connectorSettings, isFound := cache.Cache.Get(fmt.Sprintf("connectorEvse%dId%dConfiguration", connector.EvseId, connector.ConnectorId))
		if !isFound {
			continue
		}
		cachedConnector := connectorSettings.(*settings.Connector)
		if cachedConnector != nil {
			handler.notifyConnectorStatus(connector, core.ChargePointStatus(cachedConnector.Status), core.NoError)
			switch core.ChargePointStatus(cachedConnector.Status) {
			case core.ChargePointStatusPreparing:
				err = handler.startCharging(cachedConnector.Session.TagId)
				if err != nil {
					handler.notifyConnectorStatus(connector, core.ChargePointStatusAvailable, core.InternalError)
					continue
				}
				break
			case core.ChargePointStatusCharging:
				err = connector.ResumeCharging(data.Session(cachedConnector.Session))
				if err != nil {
					err = handler.stopChargingConnector(connector, core.ReasonOther)
					if err != nil {
						handler.notifyConnectorStatus(connector, core.ChargePointStatusFaulted, core.InternalError)
					}
				}
				break
			}
		}
	}
}

// isTagAuthorized Check if the tag is authorized for charging. If the auth cache is enabled and if it can preauthorize from cache,
// the program will check the cache first and reauthorize with the sendAuthorizeRequest to the central system after 10 seconds.
// If cache is not enabled, it will just execute sendAuthorizeRequest and retrieve the status from the request.
func (handler *ChargePointHandler) isTagAuthorized(tagId string) bool {
	response := false
	authCacheEnabled, err := settings.GetConfigurationValue("AuthorizationCacheEnabled")
	if err != nil {
		authCacheEnabled = "false"
	}
	localPreAuthorize, err := settings.GetConfigurationValue("LocalPreAuthorize")
	if err != nil {
		localPreAuthorize = "false"
	}
	if authCacheEnabled == "true" && localPreAuthorize == "true" {
		//Check if the tag exists in cache and is valid.
		log.Println("Authorizing tag ", tagId, " with cache")
		if data.IsTagAuthorized(tagId) {
			_, err2 := scheduler.Every(10).Seconds().LimitRunsTo(1).Do(handler.sendAuthorizeRequest, tagId)
			if err2 != nil {
				log.Println(err2)
			}
			return true
		}
	}
	//If the card is not in cache or is not authorized, (re)authorize it with the central system
	log.Println("Authorizing tag with central system: ", tagId)
	tagInfo, err := handler.sendAuthorizeRequest(tagId)
	if tagInfo != nil && tagInfo.Status == types2.AuthorizationStatusAccepted {
		response = true
	}
	log.Println("Tag authorization result: ", response)
	return response
}

// sendAuthorizeRequest Send a AuthorizeRequest to the central system to get information on the tagId authorization status.
// Adds the tag to the cache if it is enabled.
func (handler *ChargePointHandler) sendAuthorizeRequest(tagId string) (*types2.IdTagInfo, error) {
	var err error
	response, err := handler.chargePoint.SendRequest(core.AuthorizeRequest{IdTag: tagId})
	authInfo := response.(*core.AuthorizeConfirmation)
	switch authInfo.IdTagInfo.Status {
	case types2.AuthorizationStatusBlocked:
	case types2.AuthorizationStatusExpired:
	case types2.AuthorizationStatusInvalid:
		err = handler.stopChargingConnectorWithTagId(tagId, core.ReasonDeAuthorized)
	}
	value, err2 := settings.GetConfigurationValue("AuthorizationCacheEnabled")
	if err2 == nil && value == "true" {
		data.AddTag(tagId, authInfo.IdTagInfo)
	}
	return authInfo.IdTagInfo, err
}

// notifyConnectorStatus Notified the central system about the connector's status and updates the LED indicator.
func (handler *ChargePointHandler) notifyConnectorStatus(connector *Connector, status core.ChargePointStatus, errorCode core.ChargePointErrorCode) {
	if connector != nil {
		request := core.StatusNotificationRequest{
			ConnectorId: connector.ConnectorId,
			Status:      status,
			ErrorCode:   errorCode,
			Timestamp:   &types2.DateTime{Time: time.Now()},
		}
		connector.SetStatus(status)
		callback := func(confirmation ocpp.Response, protoError error) {
			log.Printf("Changed status of the connector %d to %s", connector.ConnectorId, connector.ConnectorStatus)
			connectorIndex := connector.ConnectorId - 1
			handler.displayLEDStatus(connectorIndex, connector.ConnectorStatus)
		}
		err := handler.chargePoint.SendRequestAsync(request, callback)
		if err != nil {
			log.Println("Cannot change status of connector: ", err)
			return
		}
	}
}

// HandleChargingRequest Entry point for determining if the request is to start or stop charging. Trying to find a connector that has the tag stored in the Session; if such a connector exists,
// execute stopChargingConnector, otherwise startCharging.
func (handler *ChargePointHandler) HandleChargingRequest(tagId string) {
	log.Printf("Handling request for tag %s", tagId)
	var connector = handler.FindConnectorWithTagId(tagId)
	if connector != nil {
		err := handler.stopChargingConnector(connector, core.ReasonLocal)
		if err != nil {
			log.Printf("Error stopping charing the connector: %s", err)
			return
		}
	} else {
		err := handler.startCharging(tagId)
		if err != nil {
			log.Printf("Error started charing the connector: %s", err)
			return
		}
	}
}

// startCharging Start charging on the first available connector. If there is no available connector, reject the request.
func (handler *ChargePointHandler) startCharging(tagId string) error {
	var connector = handler.FindAvailableConnector()
	if connector != nil {
		return handler.startChargingConnector(connector, tagId)
	}
	return errors.New("no available connectors")
}

// startChargingConnector Start charging a connector with the specified ID. Send the request to the central system, turn on the Connector,
// update the status of the Connector, start the timer and sample the PowerMeter if enabled.
func (handler *ChargePointHandler) startChargingConnector(connector *Connector, tagId string) error {
	if connector != nil && connector.IsAvailable() && handler.isTagAuthorized(tagId) {
		handler.notifyConnectorStatus(connector, core.ChargePointStatusPreparing, core.NoError)
		request := core.StartTransactionRequest{
			ConnectorId: connector.ConnectorId,
			IdTag:       tagId,
			Timestamp:   &types2.DateTime{Time: time.Now()},
			MeterStart:  0,
		}
		var callback = func(confirmation ocpp.Response, protoError error) {
			startTransactionConf := confirmation.(*core.StartTransactionConfirmation)
			if startTransactionConf.TransactionId > 0 && startTransactionConf.IdTagInfo.Status == types2.AuthorizationStatusAccepted {
				err := connector.StartCharging(strconv.Itoa(startTransactionConf.TransactionId), tagId)
				settings.UpdateConnectorSessionInfo(
					connector.EvseId,
					connector.ConnectorId,
					&settings.Session{
						IsActive:      connector.session.IsActive,
						TagId:         connector.session.TagId,
						TransactionId: connector.session.TransactionId,
						Started:       connector.session.Started,
						Consumption:   connector.session.Consumption,
					})
				fmt.Println(connector.session)
				if err != nil {
					log.Printf("Unable to start charging connector %d: %s", connector.ConnectorId, err)
					return
				}
				handler.notifyConnectorStatus(connector, core.ChargePointStatusCharging, core.NoError)
				log.Printf("Started charging connector %d at %s", connector.ConnectorId, time.Now())
				if connector.PowerMeterEnabled {
					measurandsString, err := settings.GetConfigurationValue("MeterValuesSampledData")
					var measurands []types2.Measurand
					if err != nil {
						measurandsString = string(types2.MeasurandPowerActiveExport)
					}
					for _, s := range strings.Split(measurandsString, ",") {
						measurands = append(measurands, types2.Measurand(s))
					}
					sampleInterval, err := settings.GetConfigurationValue("MeterValueSampleInterval")
					if err != nil {
						sampleInterval = "10"
					}
					_, err = scheduler.Every(fmt.Sprintf("%ss", sampleInterval)).
						Tag(fmt.Sprintf("connector%dSampling", connector.ConnectorId)).Do(connector.SamplePowerMeter, measurands)
				}
				_, err = scheduler.Every(connector.MaxChargingTime).Minutes().LimitRunsTo(1).
					Tag(fmt.Sprintf("connector%dTimer", connector.ConnectorId)).Do(handler.stopChargingConnector, connector, core.ReasonOther)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				log.Printf("Transaction unauthorized at connector %d", connector.ConnectorId)
			}
		}
		return handler.chargePoint.SendRequestAsync(request, callback)
	}
	return errors.New("connector unavailable or card unauthorized")
}

// stopChargingConnector Stop charging a connector with the specified ID. Update the status(es), turn off the Connector and calculate the energy consumed.
func (handler *ChargePointHandler) stopChargingConnector(connector *Connector, reason core.Reason) error {
	if connector != nil && (connector.IsCharging() || connector.IsPreparing()) {
		stopTransactionOnEVDisconnect, err := settings.GetConfigurationValue("StopTransactionOnEVSideDisconnect")
		if err != nil {
			return err
		}
		if stopTransactionOnEVDisconnect != "true" && reason == core.ReasonEVDisconnected {
			handler.notifyConnectorStatus(connector, core.ChargePointStatusSuspendedEVSE, core.NoError)
			err := connector.StopCharging()
			handler.notifyConnectorStatus(connector, core.ChargePointStatusFinishing, core.NoError)
			settings.UpdateConnectorSessionInfo(
				connector.EvseId,
				connector.ConnectorId,
				&settings.Session{
					IsActive:      connector.session.IsActive,
					TagId:         connector.session.TagId,
					TransactionId: connector.session.TransactionId,
					Started:       connector.session.Started,
					Consumption:   connector.session.Consumption,
				})
			handler.notifyConnectorStatus(connector, core.ChargePointStatusAvailable, core.NoError)
			return err
		}
		transactionId, err := strconv.Atoi(connector.session.TransactionId)
		if err != nil {
			return err
		}
		request := core.StopTransactionRequest{
			TransactionId: transactionId,
			IdTag:         "",
			MeterStop:     0,
			Timestamp:     &types2.DateTime{Time: time.Now()},
			Reason:        reason,
		}
		var callback = func(confirmation ocpp.Response, protoError error) {
			if protoError != nil {
				log.Printf("Server responded with error for stopping a transaction at %d: %s", connector.ConnectorId, err)
				return
			}
			log.Println("Stopping transaction at ", connector.ConnectorId)
			err := connector.StopCharging()
			settings.UpdateConnectorSessionInfo(
				connector.EvseId,
				connector.ConnectorId,
				&settings.Session{
					IsActive:      connector.session.IsActive,
					TagId:         connector.session.TagId,
					TransactionId: connector.session.TransactionId,
					Started:       connector.session.Started,
					Consumption:   connector.session.Consumption,
				})
			if err != nil {
				log.Printf("Unable to stop charging connector %d: %s", connector.ConnectorId, err)
				handler.notifyConnectorStatus(connector, core.ChargePointStatusFinishing, core.InternalError)
				return
			}
			_, err = scheduler.Every(1).Seconds().LimitRunsTo(1).Do(handler.notifyConnectorStatus, connector, core.ChargePointStatusFinishing, core.NoError)
			_, err = scheduler.Every(3).Seconds().LimitRunsTo(1).Do(handler.notifyConnectorStatus, connector, core.ChargePointStatusAvailable, core.NoError)
			if err != nil {
				fmt.Println(err)
			}
			_ = scheduler.RemoveByTag(fmt.Sprintf("connector%dSampling", connector.ConnectorId))
			err = scheduler.RemoveByTag(fmt.Sprintf("connector%dTimer", connector.ConnectorId))
			log.Printf("Stopped charging connector %d at %s", connector.ConnectorId, time.Now())
			//log.Printf("Transaction unauthorized at connector %d", connector.ConnectorId)
		}
		return handler.chargePoint.SendRequestAsync(request, callback)
	}
	return errors.New("connector not charging")
}

// stopChargingConnectorWithTagId Search for a Connector that contains the tagId and stop the charging.
func (handler *ChargePointHandler) stopChargingConnectorWithTagId(tagId string, reason core.Reason) error {
	var connector = handler.FindConnectorWithTagId(tagId)
	if connector != nil {
		return handler.stopChargingConnector(connector, reason)
	}
	return errors.New("no connector with tag id")
}

// stopChargingConnectorWithTransactionId Search for a Connector that contains the transactionId and stop the charging.
func (handler *ChargePointHandler) stopChargingConnectorWithTransactionId(transactionId string) error {
	var connector = handler.FindConnectorWithTransactionId(transactionId)
	if connector != nil {
		return handler.stopChargingConnector(connector, core.ReasonRemote)
	}
	return errors.New("no connector with transaction id")
}

// AddConnectors Add the Connectors from the connectors.json file to the handler. Create and add all their components and initialize the struct.
func (handler *ChargePointHandler) AddConnectors(connectors []*settings.Connector) {
	log.Println("Adding connectors")
	handler.Connectors = []*Connector{}
	for _, connector := range connectors {
		var powerMeter *hardware.PowerMeter
		powerMeterInstance, err := hardware.NewPowerMeter(
			connector.PowerMeter.PowerMeterPin,
			connector.PowerMeter.SpiBus,
			connector.PowerMeter.ShuntOffset,
			connector.PowerMeter.VoltageDividerOffset,
		)
		if err != nil {
			log.Printf("Cannot instantiate power meter: %s", err)
		} else {
			powerMeter = powerMeterInstance
		}
		connectorObj, err := NewConnector(
			connector.EvseId,
			connector.ConnectorId,
			connector.Type,
			hardware.NewRelay(connector.Relay.RelayPin, connector.Relay.InverseLogic),
			powerMeter,
			connector.PowerMeter.Enabled,
			handler.Settings.ChargePoint.Info.MaxChargingTime,
		)
		if err != nil {
			continue
		}
		handler.Connectors = append(handler.Connectors, &connectorObj)
		fmt.Print("Added connector ")
		pretty.Print(connectorObj)
	}
}

func (handler *ChargePointHandler) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	var response core.ClearCacheStatus = core.ClearCacheStatusRejected
	log.Printf("Requested clear cache")
	authCacheEnabled, err := settings.GetConfigurationValue("AuthorizationCacheEnabled")
	if err != nil {
		log.Printf("Error clearing cache: %s", err)
	} else if authCacheEnabled == "true" {
		data.RemoveCachedTags()
		response = core.ClearCacheStatusAccepted
	}
	return core.NewClearCacheConfirmation(response), err
}

func (handler *ChargePointHandler) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	var response core.DataTransferStatus = core.DataTransferStatusRejected
	return core.NewDataTransferConfirmation(response), errors.New("unsupported action")
}

func (handler *ChargePointHandler) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	log.Printf("Requested configuration")
	var unknownKeys []string
	var configArray []core.ConfigurationKey
	configuration, err := settings.GetConfiguration()
	if err == nil && configuration != nil {
		configArray = configuration.GetConfig()
		for _, key := range request.Key {
			_, err := settings.GetConfigurationValue(key)
			if err != nil {
				unknownKeys = append(unknownKeys, key)
			}
		}
	}
	return &core.GetConfigurationConfirmation{ConfigurationKey: configArray, UnknownKey: unknownKeys}, err
}

func (handler *ChargePointHandler) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {
	var response types2.RemoteStartStopStatus = types2.RemoteStartStopStatusRejected
	log.Printf("Got remote stop request for connector %d with tag %s", request.ConnectorId, request.IdTag)
	var connector = handler.FindConnectorWithId(*request.ConnectorId)
	if connector != nil {
		//Delay the charging by 3 seconds
		_, err := scheduler.Every(3).Seconds().LimitRunsTo(1).Do(handler.startCharging, request.IdTag)
		if err != nil {
			return core.NewRemoteStartTransactionConfirmation(response), err
		}
		response = types2.RemoteStartStopStatusAccepted
	}
	return core.NewRemoteStartTransactionConfirmation(response), err
}

func (handler *ChargePointHandler) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	var response types2.RemoteStartStopStatus = types2.RemoteStartStopStatusRejected
	var transactionId = string(rune(request.TransactionId))
	var connector = handler.FindConnectorWithTransactionId(transactionId)
	log.Printf("Got remote stop request for transaction %s", transactionId)
	if connector != nil {
		//Delay stopping the transaction by 3 seconds
		_, err := scheduler.Every(3).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnectorWithTransactionId, transactionId)
		if err != nil {
			return core.NewRemoteStopTransactionConfirmation(response), err
		}
		response = types2.RemoteStartStopStatusAccepted
	}
	return core.NewRemoteStopTransactionConfirmation(response), nil
}

func (handler *ChargePointHandler) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	var response core.ResetStatus = core.ResetStatusRejected
	log.Printf("Requested reset %s", request.Type)
	if request.Type == core.ResetTypeHard {
		handler.CleanUp(core.ReasonHardReset)
		_, err = scheduler.Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		if err == nil {
			response = core.ResetStatusAccepted
		}
	} else if request.Type == core.ResetTypeSoft {
		handler.CleanUp(core.ReasonSoftReset)
		_, err = scheduler.Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		if err == nil {
			response = core.ResetStatusAccepted
		}
	}
	return core.NewResetConfirmation(response), err
}

func (handler *ChargePointHandler) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	var response core.UnlockStatus = core.UnlockStatusNotSupported
	connector := handler.FindConnectorWithId(request.ConnectorId)
	_, err = scheduler.Every(1).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnector, connector.ConnectorId)
	if err == nil {
		response = core.UnlockStatusUnlocked
	}
	return core.NewUnlockConnectorConfirmation(response), err
}

func (handler *ChargePointHandler) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	var response core.AvailabilityStatus = core.AvailabilityStatusRejected
	return core.NewChangeAvailabilityConfirmation(response), nil
}

func (handler *ChargePointHandler) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	response := core.ConfigurationStatusRejected
	log.Printf("Requested configuration change")
	err = settings.UpdateKey(request.Key, request.Value)
	if err == nil {
		response = core.ConfigurationStatusAccepted
	}
	return core.NewChangeConfigurationConfirmation(response), err
}

func (handler *ChargePointHandler) OnTriggerMessage(request *remotetrigger.TriggerMessageRequest) (confirmation *remotetrigger.TriggerMessageConfirmation, err error) {
	log.Printf("received %s for %v", request.GetFeatureName(), request.RequestedMessage)
	status := remotetrigger.TriggerMessageStatusRejected
	switch request.RequestedMessage {
	case core.BootNotificationFeatureName:
		_, err = scheduler.Every(10).Seconds().LimitRunsTo(1).Do(handler.bootNotification)
		if err != nil {
			return
		}
		status = remotetrigger.TriggerMessageStatusAccepted
		break
	case firmware.DiagnosticsStatusNotificationFeatureName:
		break
	case firmware.FirmwareStatusNotificationFeatureName:
		break
	case core.HeartbeatFeatureName:
		_, err = scheduler.Every(5).Seconds().LimitRunsTo(1).Do(handler.sendHeartBeat)
		if err != nil {
			return
		}
		status = remotetrigger.TriggerMessageStatusAccepted
	case core.MeterValuesFeatureName:
		break
	case core.StatusNotificationFeatureName:
		connectorID := *request.ConnectorId
		connector := handler.FindConnectorWithId(connectorID)
		if connector != nil {
			_, err = scheduler.Every(5).Seconds().LimitRunsTo(1).Do(handler.notifyConnectorStatus, connector.ConnectorId, connector.ConnectorStatus, core.NoError)
			if err != nil {
				return nil, err
			}
		}
		status = remotetrigger.TriggerMessageStatusAccepted
		break
	default:
		return remotetrigger.NewTriggerMessageConfirmation(remotetrigger.TriggerMessageStatusNotImplemented), nil
	}
	return remotetrigger.NewTriggerMessageConfirmation(status), nil
}

func (handler *ChargePointHandler) OnReserveNow(request *reservation.ReserveNowRequest) (confirmation *reservation.ReserveNowConfirmation, err error) {
	connector := handler.FindConnectorWithId(request.ConnectorId)
	if connector == nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusUnavailable), nil
	} else if !connector.IsAvailable() {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusOccupied), nil
	}
	err = connector.ReserveConnector(request.ReservationId)
	if err != nil {
		return reservation.NewReserveNowConfirmation(reservation.ReservationStatusRejected), err
	}
	_, err = scheduler.At(request.ExpiryDate.Format("HH:MM")).Do(connector.RemoveReservation)
	return reservation.NewReserveNowConfirmation(reservation.ReservationStatusAccepted), err
}

func (handler *ChargePointHandler) OnCancelReservation(request *reservation.CancelReservationRequest) (confirmation *reservation.CancelReservationConfirmation, err error) {
	connector := handler.FindConnectorWithReservationId(request.ReservationId)
	if connector != nil {
		err = connector.RemoveReservation()
		if err != nil {
			return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), err
		}
		return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusAccepted), nil
	}
	return reservation.NewCancelReservationConfirmation(reservation.CancelReservationStatusRejected), nil
}

// bootNotification Notify the central system that the charging point is online. Set the heartbeat interval and call startup.
// If the central system does not accept the charge point, exit the client.
func (handler *ChargePointHandler) bootNotification() {
	callback := func(confirmation ocpp.Response, protoError error) {
		bootConf := confirmation.(*core.BootNotificationConfirmation)
		if bootConf.Status == core.RegistrationStatusAccepted {
			log.Printf("Notified and accepted from the central system")
			heartBeatInterval, _ := settings.GetConfigurationValue("HeartbeatInterval")
			if bootConf.Interval > 0 {
				heartBeatInterval = fmt.Sprintf("%d", bootConf.Interval)
			}
			heartBeatInterval = fmt.Sprintf("%ss", heartBeatInterval)
			_, err := scheduler.Every(heartBeatInterval).Do(handler.sendHeartBeat)
			if err != nil {
				fmt.Println(err)
			}
			handler.startup()
		} else {
			err := RetrySendingRequest("BootNotificationRetries", -1, bootConf.Interval, handler.bootNotification, nil)
			if err == nil {
				return
			}
			log.Printf("Denied by the central system.")
			os.Exit(-1)
		}
	}
	request := core.BootNotificationRequest{
		ChargePointModel:  handler.Settings.ChargePoint.Info.Model,
		ChargePointVendor: handler.Settings.ChargePoint.Info.Vendor,
	}
	err := handler.chargePoint.SendRequestAsync(request, callback)
	if err != nil {
		return
	}
}

func RetrySendingRequest(cacheRetryKey string, maxRetries int, interval int, function interface{}, functionParams ...interface{}) (err error) {
	var maxMessageAttempts string
	var retryInterval interface{}
	retries, isFound := cache.Cache.Get(cacheRetryKey)
	if !isFound {
		cache.Cache.Set(cacheRetryKey, 0, goCache.DefaultExpiration)
	}
	err = cache.Cache.Increment(cacheRetryKey, 1)
	if err != nil {
		return err
	}
	if maxRetries < 0 {
		maxMessageAttempts, err = settings.GetConfigurationValue("TransactionMessageAttempts")
		maxRetries, err = strconv.Atoi(maxMessageAttempts)
		if err != nil {
			maxRetries = 5
		}
	}
	if retries.(int) < maxRetries {
		retryInterval, err = settings.GetConfigurationValue("TransactionMessageRetryInterval")
		if err != nil {
			retryInterval = "30"
		}
		if interval > 0 {
			retryInterval = fmt.Sprintf("%d", interval)
		}
		retryInterval = fmt.Sprintf("%ss", retryInterval)
		_, err = scheduler.Every(retryInterval).LimitRunsTo(1).Tag(fmt.Sprintf("%s-#%d", cacheRetryKey, retries)).Do(function, functionParams)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("max retries reached")
}

// sendHeartBeat Send a heartbeat to the central system.
func (handler *ChargePointHandler) sendHeartBeat() error {
	err := handler.chargePoint.SendRequestAsync(core.HeartbeatRequest{},
		func(confirmation ocpp.Response, protoError error) {
			log.Printf("Sent heartbeat")
		})
	return err
}

// CleanUp When exiting the client, stop all the transactions, clean up all the peripherals and terminate the connection.
func (handler *ChargePointHandler) CleanUp(reason core.Reason) {
	log.Println("Cleaning up ChargePoint, reason:", reason)
	for _, connector := range handler.Connectors {
		if connector.IsCharging() {
			log.Println("Stopping a transaction at connector: ", connector.ConnectorId)
			err := handler.stopChargingConnector(connector, reason)
			if err != nil {
				log.Printf("error while stopping the transaction at cleanup: %v", err)
			}
		}
	}
	log.Println("Clearing the scheduler...")
	scheduler.Stop()
	scheduler.Clear()
	log.Println("Disconnecting the client..")
	handler.chargePoint.Stop()
	if handler.TagReader != nil {
		log.Println("Cleaning up Tag Reader")
		close(handler.TagReader.TagChannel)
		handler.TagReader.Cleanup()
	}
	if handler.LCD != nil {
		log.Println("Clearing LCD")
		close(handler.LCD.LCDChannel)
		handler.LCD.Cleanup()
	}
	if handler.LEDStrip != nil {
		log.Println("Clearing LEDs")
		handler.LEDStrip.Cleanup()
	}
}

func GetTLSClient(CACertificatePath string, ClientCertificatePath string, ClientKeyPath string) *ws.Client {
	certPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	// Load CA cert
	caCert, err := ioutil.ReadFile(CACertificatePath)
	if err != nil {
		log.Println(err)
		return nil
	} else if !certPool.AppendCertsFromPEM(caCert) {
		log.Println("no ca.cert file found, will use system CA certificates")
		return nil
	}
	// Load client certificate
	certificate, err := tls.LoadX509KeyPair(ClientCertificatePath, ClientKeyPath)
	if err != nil {
		log.Printf("couldn't load client TLS certificate: %v \n", err)
		return nil
	}
	// Create client with TLS config
	return ws.NewTLSClient(&tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{certificate},
	})
}

func (handler *ChargePointHandler) Run() {
	var client ws.WsClient = nil
	if handler.Settings.ChargePoint.Info.TLS.IsEnabled {
		client = GetTLSClient(handler.Settings.ChargePoint.Info.TLS.CACertificatePath, handler.Settings.ChargePoint.Info.TLS.ClientCertificatePath, handler.Settings.ChargePoint.Info.TLS.ClientKeyPath)
		handler.chargePoint = ocpp16.NewChargePoint(handler.Settings.ChargePoint.Info.Id, nil, client)
	} else {
		handler.chargePoint = ocpp16.NewChargePoint(handler.Settings.ChargePoint.Info.Id, nil, nil)
	}
	handler.chargePoint.SetCoreHandler(handler)
	handler.chargePoint.SetReservationHandler(handler)
	handler.chargePoint.SetRemoteTriggerHandler(handler)
	maCachedTagsString, err := settings.GetConfigurationValue("MaxCachedTags")
	maxCachedTags, err := strconv.Atoi(maCachedTagsString)
	if err == nil {
		data.SetMaxCachedTags(maxCachedTags)
	}
	serverUrl := fmt.Sprintf("ws://%s/%s", handler.Settings.ChargePoint.Info.ServerUri, handler.Settings.ChargePoint.Info.Id)
	log.Println("Trying to connect to the central system: ", serverUrl)
	connectErr := handler.chargePoint.Start(serverUrl)
	go handler.listenForTag()
	if connectErr != nil {
		log.Printf("Error connecting to the central system: %s", connectErr)
		handler.CleanUp(core.ReasonOther)
		handler.chargePoint.Stop()
	} else {
		log.Printf("connected to central server: %s \n with ID: %s", serverUrl, handler.Settings.ChargePoint.Info.Id)
		handler.IsAvailable = true
		handler.bootNotification()
	}
	/*	time.Sleep(time.Second * 10)
		handler.HandleChargingRequest("F5ED1377")
		time.Sleep(time.Second * 10)
		handler.HandleChargingRequest("F5ED1377")*/
}

// listenForTag Listen for a RFID/NFC tag on a separate thread. If a tag is detected, call the HandleChargingRequest.
// Blink the LED if indication is enabled.
func (handler *ChargePointHandler) listenForTag() {
	if !handler.Settings.ChargePoint.Hardware.TagReader.IsSupported {
		return
	}
	for {
		fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
		select {
		case tagId := <-handler.TagReader.TagChannel:
			handler.indicateCard(len(handler.Connectors), hardware.WHITE)
			handler.HandleChargingRequest(strings.ToUpper(tagId))
			continue
		default:
			time.Sleep(time.Millisecond * 300)
			continue
		}
	}
}

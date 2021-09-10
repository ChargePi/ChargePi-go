package chargepoint

import (
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
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type ChargePointHandler struct {
	chargePoint      ocpp16.ChargePoint
	IsAvailable      bool
	Connectors       []*Connector
	Settings         *settings.Settings
	TagReader        *hardware.TagReader
	LEDStrip         *hardware.LEDStrip
	LCD              *hardware.LCD
	connectorChannel chan rxgo.Item
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

// startCharging Start charging on the first available Connector. If there is no available Connector, reject the request.
func (handler *ChargePointHandler) startCharging(tagId string) error {
	var connector = handler.FindAvailableConnector()
	if connector != nil {
		return handler.startChargingConnector(connector, tagId)
	}
	return errors.New("no available connectors")
}

// startChargingConnector Start charging a connector with the specified ID. Send the request to the central system, turn on the Connector,
// update the status of the Connector, and start the maxChargingTime timer and sample the PowerMeter, if it's enabled.
func (handler *ChargePointHandler) startChargingConnector(connector *Connector, tagId string) error {
	if connector == nil {
		return fmt.Errorf("connector is nil")
	}
	if !connector.IsAvailable() || !handler.IsAvailable {
		return errors.New("connector or cp unavailable")
	}
	if handler.isTagAuthorized(tagId) {
		request := core.StartTransactionRequest{
			ConnectorId: connector.ConnectorId,
			IdTag:       tagId,
			Timestamp:   &types2.DateTime{Time: time.Now()},
			MeterStart:  0,
		}
		callback := func(confirmation ocpp.Response, protoError error) {
			startTransactionConf := confirmation.(*core.StartTransactionConfirmation)
			if startTransactionConf.TransactionId > 0 && startTransactionConf.IdTagInfo.Status == types2.AuthorizationStatusAccepted {
				err := connector.StartCharging(strconv.Itoa(startTransactionConf.TransactionId), tagId)
				if err != nil {
					log.Printf("Unable to start charging connector %d: %s", connector.ConnectorId, err)
					return
				}
				log.Printf("Started charging connector %d at %s", connector.ConnectorId, time.Now())
				if connector.PowerMeterEnabled {
					sampleError := preparePowerMeterAtConnector(connector)
					if sampleError != nil {
						log.Printf("Cannot sample connector %d; %v \n", connector.ConnectorId, err)
					}
				}
				// schedule timer to stop the transaction at the time limit
				_, err = scheduler.Every(connector.MaxChargingTime).Minutes().LimitRunsTo(1).
					Tag(fmt.Sprintf("connector%dTimer", connector.ConnectorId)).Do(handler.stopChargingConnector, connector, core.ReasonOther)
				if err != nil {
					fmt.Println("cannot schedule stop charging:", err)
				}
			} else {
				log.Printf("Transaction unauthorized at connector %d", connector.ConnectorId)
			}
		}
		return handler.chargePoint.SendRequestAsync(request, callback)
	}
	return errors.New("card unauthorized")
}

// stopChargingConnector Stop charging a connector with the specified ID. Update the status(es), turn off the Connector and calculate the energy consumed.
func (handler *ChargePointHandler) stopChargingConnector(connector *Connector, reason core.Reason) error {
	if connector == nil {
		return fmt.Errorf("connector pointer is nil")
	}
	if connector.IsCharging() || connector.IsPreparing() {
		stopTransactionOnEVDisconnect, err := settings.GetConfigurationValue("StopTransactionOnEVSideDisconnect")
		if err != nil {
			return err
		}
		if stopTransactionOnEVDisconnect != "true" && reason == core.ReasonEVDisconnected {
			err := connector.StopCharging(reason)
			return err
		}
		transactionId, err := strconv.Atoi(connector.session.TransactionId)
		if err != nil {
			return err
		}
		request := core.StopTransactionRequest{
			TransactionId: transactionId,
			IdTag:         "",
			MeterStop:     int(connector.session.CalculateEnergyConsumptionWithAvgPower()),
			Timestamp:     &types2.DateTime{Time: time.Now()},
			Reason:        reason,
		}
		var callback = func(confirmation ocpp.Response, protoError error) {
			if protoError != nil {
				log.Printf("Server responded with error for stopping a transaction at %d: %s", connector.ConnectorId, err)
				return
			}
			log.Println("Stopping transaction at ", connector.ConnectorId)
			err := connector.StopCharging(reason)
			if err != nil {
				log.Printf("Unable to stop charging connector %d: %s", connector.ConnectorId, err)
				return
			}
			_ = scheduler.RemoveByTag(fmt.Sprintf("connector%dSampling", connector.ConnectorId))
			err = scheduler.RemoveByTag(fmt.Sprintf("connector%dTimer", connector.ConnectorId))
			log.Printf("Stopped charging connector %d at %s", connector.ConnectorId, time.Now())
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
		handler.Connectors = append(handler.Connectors, connectorObj)
		fmt.Print("Added connector ")
		pretty.Print(connectorObj)
		// turn off the relay at boot
		connectorObj.relay.Off()
	}
}

func (handler *ChargePointHandler) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	var (
		response         = core.ClearCacheStatusRejected
		authCacheEnabled string
	)
	log.Printf("Requested clear cache")
	authCacheEnabled, err = settings.GetConfigurationValue("AuthorizationCacheEnabled")
	if err != nil {
		log.Printf("Error clearing cache: %s", err)
	} else if authCacheEnabled == "true" {
		data.RemoveCachedTags()
		response = core.ClearCacheStatusAccepted
	}
	return core.NewClearCacheConfirmation(response), err
}

func (handler *ChargePointHandler) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	var response = core.DataTransferStatusRejected
	return core.NewDataTransferConfirmation(response), errors.New("unsupported action")
}

func (handler *ChargePointHandler) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	var (
		unknownKeys   []string
		configArray   []core.ConfigurationKey
		configuration *settings.OCPPConfig
	)
	log.Printf("Requested configuration")
	configuration, err = settings.GetConfiguration()
	if err == nil && configuration != nil {
		configArray = configuration.GetConfig()
		for _, key := range request.Key {
			_, keyErr := settings.GetConfigurationValue(key)
			if keyErr != nil {
				unknownKeys = append(unknownKeys, key)
			}
		}
	}
	return &core.GetConfigurationConfirmation{ConfigurationKey: configArray, UnknownKey: unknownKeys}, err
}

func (handler *ChargePointHandler) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {
	var (
		connectorId = *request.ConnectorId
		response    = types2.RemoteStartStopStatusRejected
		connector   = handler.FindConnectorWithId(connectorId)
	)
	log.Printf("Got remote start request for connector %d with tag %s", connectorId, request.IdTag)
	if connector != nil {
		//Delay the charging by 3 seconds
		_, err = scheduler.Every(3).Seconds().LimitRunsTo(1).Do(handler.startCharging, request.IdTag)
		if err != nil {
			return core.NewRemoteStartTransactionConfirmation(response), err
		}
		response = types2.RemoteStartStopStatusAccepted
	}
	return core.NewRemoteStartTransactionConfirmation(response), err
}

func (handler *ChargePointHandler) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	var (
		response      = types2.RemoteStartStopStatusRejected
		transactionId = fmt.Sprintf("%d", request.TransactionId)
		connector     = handler.FindConnectorWithTransactionId(transactionId)
	)
	log.Printf("Got remote stop request for transaction %s", transactionId)
	if connector != nil {
		//Delay stopping the transaction by 3 seconds
		_, err = scheduler.Every(3).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnectorWithTransactionId, transactionId)
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
		_, err = scheduler.Every(3).Seconds().LimitRunsTo(1).Do(handler.CleanUp, core.ReasonHardReset)
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
	if connector != nil {
		_, err = scheduler.Every(1).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnector, connector, core.ReasonUnlockCommand)
		if err == nil {
			response = core.UnlockStatusUnlocked
		}
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
		_, err = scheduler.Every(5).Seconds().LimitRunsTo(1).Do(handler.bootNotification)
		if err != nil {
			break
		}
		status = remotetrigger.TriggerMessageStatusAccepted
		break
	case firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
		break
	case core.HeartbeatFeatureName:
		_, err = scheduler.Every(5).Seconds().LimitRunsTo(1).Do(handler.sendHeartBeat)
		if err != nil {
			break
		}
		status = remotetrigger.TriggerMessageStatusAccepted
		break
	case core.MeterValuesFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
		break
	case core.StatusNotificationFeatureName:
		if request.ConnectorId == nil {
			// send the status of all connectors after the response
			defer func() {
				for _, connector := range handler.Connectors {
					if connector.connectorNotificationChannel != nil {
						connector.connectorNotificationChannel <- rxgo.Of(connector)
					}
				}
			}()
			status = remotetrigger.TriggerMessageStatusAccepted
			return remotetrigger.NewTriggerMessageConfirmation(status), nil
		}
		connectorID := *request.ConnectorId
		connector := handler.FindConnectorWithId(connectorID)
		if connector != nil {
			defer func() {
				handler.notifyConnectorStatus(connector)
			}()
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

// bootNotification Notify the central system that the charging point is online. Set the heartbeat interval and call restoreState.
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
			handler.restoreState()
		} else {
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
	close(handler.connectorChannel)
	log.Println("Clearing the scheduler...")
	scheduler.Stop()
	scheduler.Clear()
}

// connect to the central system and attempt to boot
func (handler *ChargePointHandler) connect() {
	serverUrl := fmt.Sprintf("ws://%s/%s", handler.Settings.ChargePoint.Info.ServerUri, handler.Settings.ChargePoint.Info.Id)
	log.Println("Trying to connect to the central system: ", serverUrl)
	connectErr := handler.chargePoint.Start(serverUrl)

	go handler.listenForTag()
	handler.connectorChannel = make(chan rxgo.Item)

	// Check if the connection was successful
	if connectErr != nil {
		log.Printf("Error connecting to the central system: %s", connectErr)
		handler.CleanUp(core.ReasonOther)
		handler.chargePoint.Stop()
	} else {
		log.Printf("connected to central server: %s \n with ID: %s", serverUrl, handler.Settings.ChargePoint.Info.Id)
		handler.IsAvailable = true
		go handler.listenForConnectorStatusChange(rxgo.FromChannel(handler.connectorChannel))
		handler.bootNotification()
	}
}

func (handler *ChargePointHandler) Run() {
	// Check if the client has TLS
	var client ws.WsClient = nil
	if handler.Settings.ChargePoint.Info.TLS.IsEnabled {
		client = GetTLSClient(handler.Settings.ChargePoint.Info.TLS.CACertificatePath, handler.Settings.ChargePoint.Info.TLS.ClientCertificatePath, handler.Settings.ChargePoint.Info.TLS.ClientKeyPath)
		handler.chargePoint = ocpp16.NewChargePoint(handler.Settings.ChargePoint.Info.Id, nil, client)
	} else {
		handler.chargePoint = ocpp16.NewChargePoint(handler.Settings.ChargePoint.Info.Id, nil, nil)
	}
	// Set handlers for Core, Reservation and RemoteTrigger
	handler.chargePoint.SetCoreHandler(handler)
	handler.chargePoint.SetReservationHandler(handler)
	handler.chargePoint.SetRemoteTriggerHandler(handler)

	handler.setMaxCachedTags()
	handler.connect()
}

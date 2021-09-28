package chargepoint

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/firmware"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/reservation"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/reactivex/rxgo/v2"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"os/exec"
)

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
	if connector != nil && connector.IsAvailable() {
		//Delay the charging by 3 seconds
		_, err = scheduler.GetScheduler().Every(3).Seconds().LimitRunsTo(1).Do(handler.startCharging, request.IdTag)
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
	if connector != nil && connector.IsCharging() {
		//Delay stopping the transaction by 3 seconds
		_, err = scheduler.GetScheduler().Every(3).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnectorWithTransactionId, transactionId)
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
		_, err = scheduler.GetScheduler().Every(3).Seconds().LimitRunsTo(1).Do(handler.CleanUp, core.ReasonHardReset)
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		if err == nil {
			response = core.ResetStatusAccepted
		}
	} else if request.Type == core.ResetTypeSoft {
		handler.CleanUp(core.ReasonSoftReset)
		//todo restart ChargePi only
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
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
		_, err = scheduler.GetScheduler().Every(1).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnector, connector, core.ReasonUnlockCommand)
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
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(handler.bootNotification)
		if err != nil {
			break
		}
		status = remotetrigger.TriggerMessageStatusAccepted
		break
	case firmware.DiagnosticsStatusNotificationFeatureName, firmware.FirmwareStatusNotificationFeatureName:
		status = remotetrigger.TriggerMessageStatusNotImplemented
		break
	case core.HeartbeatFeatureName:
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(handler.sendHeartBeat)
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
	_, err = scheduler.GetScheduler().At(request.ExpiryDate.Format("HH:MM")).Do(connector.RemoveReservation)
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

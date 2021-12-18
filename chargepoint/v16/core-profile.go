package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/conf-manager"
	"github.com/xBlaz3kx/ChargePi-go/data/auth"
	"log"
	"os/exec"
)

func (handler *ChargePointHandler) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	var response core.AvailabilityStatus = core.AvailabilityStatusRejected
	return core.NewChangeAvailabilityConfirmation(response), nil
}

func (handler *ChargePointHandler) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	response := core.ConfigurationStatusRejected
	log.Printf("Requested configuration change")

	err = conf_manager.UpdateKey(request.Key, request.Value)
	if err == nil {
		response = core.ConfigurationStatusAccepted
	}

	return core.NewChangeConfigurationConfirmation(response), nil
}

func (handler *ChargePointHandler) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	log.Printf("Got request: %s", request.GetFeatureName())
	var (
		authCacheEnabled, confErr = conf_manager.GetConfigurationValue("AuthorizationCacheEnabled")
	)

	if confErr != nil || authCacheEnabled == "false" {
		log.Printf("Error clearing cache: %s", err)
		return core.NewClearCacheConfirmation(core.ClearCacheStatusRejected), nil
	}

	if authCacheEnabled == "true" {
		auth.RemoveCachedTags()
	}

	return core.NewClearCacheConfirmation(core.ClearCacheStatusAccepted), nil
}

func (handler *ChargePointHandler) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	var response = core.DataTransferStatusRejected
	return core.NewDataTransferConfirmation(response), nil
}

func (handler *ChargePointHandler) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	log.Printf("Got request: %s", request.GetFeatureName())
	var (
		unknownKeys            []string
		configArray            []core.ConfigurationKey
		response               *core.GetConfigurationConfirmation
		configuration, confErr = conf_manager.GetConfiguration()
	)

	if confErr == nil && configuration != nil {
		configArray = configuration.GetConfig()
		for _, key := range request.Key {
			_, keyErr := conf_manager.GetConfigurationValue(key)
			if keyErr != nil {
				unknownKeys = append(unknownKeys, key)
			}
		}
	}

	response = core.NewGetConfigurationConfirmation(configArray)
	response.UnknownKey = unknownKeys

	return response, nil
}

func (handler *ChargePointHandler) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	var response = core.ResetStatusRejected
	log.Printf("Requested reset %s", request.Type)

	switch request.Type {
	case core.ResetTypeHard:
		_, err = scheduler.GetScheduler().Every(3).Seconds().LimitRunsTo(1).Do(handler.CleanUp, core.ReasonHardReset)
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		if err == nil {
			response = core.ResetStatusAccepted
		}
		break
	case core.ResetTypeSoft:
		_, err = scheduler.GetScheduler().Every(3).Seconds().LimitRunsTo(1).Do(handler.CleanUp, core.ReasonSoftReset)
		// todo restart ChargePi client only
		_, err = scheduler.GetScheduler().Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		if err == nil {
			response = core.ResetStatusAccepted
		}
		break
	}

	return core.NewResetConfirmation(response), nil
}

func (handler *ChargePointHandler) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	var (
		response  = core.UnlockStatusNotSupported
		connector = handler.connectorManager.FindConnector(1, request.ConnectorId)
	)

	if connector == nil {
		return core.NewUnlockConnectorConfirmation(core.UnlockStatusUnlockFailed), nil
	}

	response = core.UnlockStatusUnlocked

	_, err = scheduler.GetScheduler().Every(1).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnector, connector, core.ReasonUnlockCommand)
	if err != nil {
		response = core.UnlockStatusUnlockFailed
	}

	return core.NewUnlockConnectorConfirmation(response), nil
}

func (handler *ChargePointHandler) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	var (
		response      = types2.RemoteStartStopStatusRejected
		transactionId = fmt.Sprintf("%d", request.TransactionId)
		connector     = handler.connectorManager.FindConnectorWithTransactionId(transactionId)
	)
	log.Printf("Got remote stop request for transaction %s", transactionId)

	if connector != nil && connector.IsCharging() {
		response = types2.RemoteStartStopStatusAccepted
		// Delay stopping the transaction by 3 seconds
		_, schedulerErr := scheduler.GetScheduler().Every(3).Seconds().LimitRunsTo(1).Do(handler.stopChargingConnectorWithTransactionId, transactionId)
		if schedulerErr != nil {
			response = types2.RemoteStartStopStatusRejected
		}
	}

	return core.NewRemoteStopTransactionConfirmation(response), nil
}

func (handler *ChargePointHandler) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {
	var (
		connectorId = *request.ConnectorId
		response    = types2.RemoteStartStopStatusRejected
		connector   = handler.connectorManager.FindConnector(1, connectorId)
	)
	log.Printf("Got remote start request for connector %d with tag %s", connectorId, request.IdTag)

	if connector != nil && connector.IsAvailable() {
		// Delay the charging by 3 seconds
		response = types2.RemoteStartStopStatusAccepted
		_, schedulerErr := scheduler.GetScheduler().Every(3).Seconds().LimitRunsTo(1).Do(handler.startCharging, request.IdTag)
		if schedulerErr != nil {
			response = types2.RemoteStartStopStatusRejected
		}
	}

	return core.NewRemoteStartTransactionConfirmation(response), nil
}

package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
	ocppManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"os/exec"
)

func (cp *ChargePoint) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	var response core.AvailabilityStatus = core.AvailabilityStatusRejected
	if err != nil {
		return core.NewChangeAvailabilityConfirmation(response), nil
	}

	if request.ConnectorId == 0 {
		// todo check if there are ongoing transactions
		cp.availability = request.Type
		response = core.AvailabilityStatusAccepted
	} else {
		// todo
	}

	return core.NewChangeAvailabilityConfirmation(response), nil
}

func (cp *ChargePoint) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	var response = core.ConfigurationStatusRejected

	log.Infof("Received request %s", request.GetFeatureName())

	err = ocppManager.UpdateKey(request.Key, request.Value)
	if err == nil {
		response = core.ConfigurationStatusAccepted
	}

	err = ocppManager.UpdateConfigurationFile()
	if err != nil {
		response = core.ConfigurationStatusRejected
	}

	return core.NewChangeConfigurationConfirmation(response), nil
}

func (cp *ChargePoint) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	log.Infof("Received request %s", request.GetFeatureName())

	var (
		response                  = core.ClearCacheStatusRejected
		authCacheEnabled, confErr = ocppManager.GetConfigurationValue(v16.AuthorizationCacheEnabled.String())
	)

	if confErr != nil || authCacheEnabled == "false" {
		log.WithError(confErr).Errorf("Cannot clear cache")
		return core.NewClearCacheConfirmation(response), nil
	}

	if authCacheEnabled == "true" {
		cp.authCache.RemoveCachedTags()
		response = core.ClearCacheStatusAccepted
	}

	return core.NewClearCacheConfirmation(core.ClearCacheStatusAccepted), nil
}

func (cp *ChargePoint) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	log.Infof("Received request %s", request.GetFeatureName())
	var response = core.DataTransferStatusRejected

	return core.NewDataTransferConfirmation(response), nil
}

func (cp *ChargePoint) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	log.Infof("Received request %s", request.GetFeatureName())

	var (
		unknownKeys            []string
		configArray            []core.ConfigurationKey
		response               *core.GetConfigurationConfirmation
		configuration, confErr = ocppManager.GetConfiguration()
	)

	if confErr == nil && configuration != nil {
		configArray = configuration
		for _, key := range request.Key {
			_, keyErr := ocppManager.GetConfigurationValue(key)
			if keyErr != nil {
				unknownKeys = append(unknownKeys, key)
			}
		}
	}

	response = core.NewGetConfigurationConfirmation(configArray)
	response.UnknownKey = unknownKeys

	return response, nil
}

func (cp *ChargePoint) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	log.Infof("Received request %s", request.GetFeatureName())
	var response = core.ResetStatusRejected

	switch request.Type {
	case core.ResetTypeHard:
		_, err = cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.CleanUp, core.ReasonHardReset)
		_, err = cp.scheduler.Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		if err == nil {
			response = core.ResetStatusAccepted
		}
		break
	case core.ResetTypeSoft:
		_, err = cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.CleanUp, core.ReasonSoftReset)
		// todo restart ChargePi client only
		_, err = cp.scheduler.Every(5).Seconds().LimitRunsTo(1).Do(exec.Command, "sudo reboot")
		if err == nil {
			response = core.ResetStatusAccepted
		}
		break
	}

	return core.NewResetConfirmation(response), nil
}

func (cp *ChargePoint) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	log.Infof("Received request %s", request.GetFeatureName())

	var (
		response = core.UnlockStatusNotSupported
		conn     = cp.connectorManager.FindConnector(1, request.ConnectorId)
	)

	if util.IsNilInterfaceOrPointer(conn) {
		return core.NewUnlockConnectorConfirmation(core.UnlockStatusUnlockFailed), nil
	}

	response = core.UnlockStatusUnlocked

	_, err = cp.scheduler.Every(1).Seconds().LimitRunsTo(1).Do(cp.stopChargingConnector, conn, core.ReasonUnlockCommand)
	if err != nil {
		response = core.UnlockStatusUnlockFailed
	}

	return core.NewUnlockConnectorConfirmation(response), nil
}

func (cp *ChargePoint) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	log.WithField("transactionId", request.TransactionId).Infof("Received request %s", request.GetFeatureName())

	var (
		response      = types.RemoteStartStopStatusRejected
		transactionId = fmt.Sprintf("%d", request.TransactionId)
		conn          = cp.connectorManager.FindConnectorWithTransactionId(transactionId)
	)

	if util.IsNilInterfaceOrPointer(conn) && conn.IsCharging() {
		response = types.RemoteStartStopStatusAccepted
		// Delay stopping the transaction by 3 seconds
		_, schedulerErr := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.stopChargingConnectorWithTransactionId, transactionId)
		if schedulerErr != nil {
			response = types.RemoteStartStopStatusRejected
		}
	}

	return core.NewRemoteStopTransactionConfirmation(response), nil
}

func (cp *ChargePoint) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {

	var (
		logInfo = log.WithFields(log.Fields{
			"connectorId": request.ConnectorId,
			"tagId":       request.IdTag,
		})
		response = types.RemoteStartStopStatusRejected
		conn     connector.Connector
	)

	logInfo.Infof("Received request %s", request.GetFeatureName())

	if request.ConnectorId != nil {
		conn = cp.connectorManager.FindConnector(1, *request.ConnectorId)
	} else {
		conn = cp.connectorManager.FindAvailableConnector()
	}

	if !util.IsNilInterfaceOrPointer(conn) && conn.IsAvailable() {
		// Delay the charging by 3 seconds
		response = types.RemoteStartStopStatusAccepted
		_, schedulerErr := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.startChargingConnector, conn, request.IdTag)
		if schedulerErr != nil {
			response = types.RemoteStartStopStatusRejected
		}
	}

	return core.NewRemoteStartTransactionConfirmation(response), nil
}

package v16

import (
	"fmt"
	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"strconv"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	ocppManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/configuration"
)

func (cp *ChargePoint) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	var response = core.AvailabilityStatusRejected

	if request.ConnectorId == 0 {
		// todo check if there are ongoing transactions
		cp.availability = request.Type
		return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusAccepted), nil
	}

	return core.NewChangeAvailabilityConfirmation(response), nil
}

func (cp *ChargePoint) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	var response = core.ConfigurationStatusRejected

	err = ocppManager.UpdateKey(request.Key, &request.Value)
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
	cp.logger.Infof("Received request %s", request.GetFeatureName())

	var (
		response                  = core.ClearCacheStatusRejected
		authCacheEnabled, confErr = ocppManager.GetConfigurationValue(v16.AuthorizationCacheEnabled.String())
	)

	if confErr != nil {
		cp.logger.WithError(confErr).Errorf("Cannot clear cache")
		return core.NewClearCacheConfirmation(response), nil
	}

	if authCacheEnabled != nil && *authCacheEnabled == "true" {
		cp.tagManager.ClearCache()
		response = core.ClearCacheStatusAccepted
	}

	return core.NewClearCacheConfirmation(response), nil
}

func (cp *ChargePoint) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	var response = core.DataTransferStatusRejected

	return core.NewDataTransferConfirmation(response), nil
}

func (cp *ChargePoint) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())

	var (
		unknownKeys            []string
		configArray            = []core.ConfigurationKey{}
		response               = core.NewGetConfigurationConfirmation(configArray)
		configuration, confErr = ocppManager.GetConfiguration()
	)

	if confErr != nil || configuration == nil {
		return response, nil
	}

	configArray = configuration

	// Get all configuration variables
	if request.Key == nil || len(request.Key) == 0 {
		response.ConfigurationKey = configArray
		response.UnknownKey = unknownKeys
		return response, nil
	}

	configArray2 := []core.ConfigurationKey{}
	// Get only the requested variables
	for _, key := range request.Key {

		// Note: redundant looping, should've just created an ocpp.ConfigurationKey function
		// Check if the key exists
		_, keyErr := ocppManager.GetConfigurationValue(key)
		if keyErr != nil {
			unknownKeys = append(unknownKeys, key)
			continue
		}

		// Key should exist, therefore find it in the config
		for _, configurationKey := range configArray {
			if key == configurationKey.Key {
				configArray2 = append(configArray2, configurationKey)
			}
		}
	}

	response.ConfigurationKey = configArray2
	response.UnknownKey = unknownKeys
	return response, nil
}

func (cp *ChargePoint) OnReset(request *core.ResetRequest) (confirmation *core.ResetConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	var response = core.ResetStatusRejected
	var retries = 3

	resetRetries, _ := ocppManager.GetConfigurationValue(v16.ResetRetries.String())
	if resetRetries != nil {
		aRetries, err := strconv.Atoi(*resetRetries)
		if err == nil {
			retries = aRetries
		}
	}

	resetErr := retry.Do(
		func() error {
			return cp.Reset(string(request.Type))
		},
		retry.Attempts(uint(retries)),
		retry.Delay(time.Second*10),
	)
	if resetErr == nil {
		response = core.ResetStatusAccepted
	}

	return core.NewResetConfirmation(response), nil
}

func (cp *ChargePoint) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())

	var (
		response   = core.UnlockStatusNotSupported
		conn, fErr = cp.connectorManager.FindEVSE(request.ConnectorId)
	)

	if fErr != nil {
		return core.NewUnlockConnectorConfirmation(response), nil
	}

	response = core.UnlockStatusUnlocked

	_, schedulerErr := cp.scheduler.Every(1).Seconds().LimitRunsTo(1).Do(cp.stopChargingConnector, conn, core.ReasonUnlockCommand)
	if schedulerErr != nil {
		response = core.UnlockStatusUnlockFailed
	}

	return core.NewUnlockConnectorConfirmation(response), nil
}

func (cp *ChargePoint) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	cp.logger.WithField("transactionId", request.TransactionId).Infof("Received request %s", request.GetFeatureName())

	var (
		response      = types.RemoteStartStopStatusRejected
		transactionId = fmt.Sprintf("%d", request.TransactionId)
		conn, fErr    = cp.connectorManager.FindEVSEWithTransactionId(transactionId)
	)

	if fErr == nil && conn.IsCharging() {
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
		logInfo = cp.logger.WithFields(log.Fields{
			"connectorId": request.ConnectorId,
			"tagId":       request.IdTag,
		})
		response = types.RemoteStartStopStatusRejected
		conn     evse.EVSE
	)

	logInfo.Infof("Received request %s", request.GetFeatureName())

	// If the connector is specified, check if it exists and is available.
	if request.ConnectorId != nil {
		conn, err = cp.connectorManager.FindEVSE(*request.ConnectorId)
	} else {
		conn, err = cp.connectorManager.FindAvailableEVSE()
	}

	if err == nil && conn.IsAvailable() {
		// Delay the charging by 3 seconds
		response = types.RemoteStartStopStatusAccepted
		_, schedulerErr := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.remoteStart, conn, 1, request.IdTag)
		if schedulerErr != nil {
			response = types.RemoteStartStopStatusRejected
		}
	}

	return core.NewRemoteStartTransactionConfirmation(response), nil
}

func (cp *ChargePoint) remoteStart(evseId, connectorId int, tagId string) {
	// todo AuthorizeRemoteTxRequests variable
	logInfo := cp.logger.WithFields(log.Fields{
		"evseId":      evseId,
		"connectorId": connectorId,
		"tagId":       tagId,
	})

	if cp.availability != core.AvailabilityTypeOperative {
		return
	}

	if !cp.isTagAuthorized(tagId) {
		return
	}

	request := core.NewStartTransactionRequest(
		evseId,
		tagId,
		0,
		types.NewDateTime(time.Now()),
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logInfo.WithError(protoError).Errorf("Central system responded with an error for %s", confirmation.GetFeatureName())
			return
		}

		startTransactionConf := confirmation.(*core.StartTransactionConfirmation)

		switch startTransactionConf.IdTagInfo.Status {
		case types.AuthorizationStatusAccepted, types.AuthorizationStatusConcurrentTx:

			// Attempt to start charging
			err := cp.connectorManager.StartCharging(evseId, tagId, strconv.Itoa(startTransactionConf.TransactionId))
			if err != nil {
				logInfo.WithError(err).Errorf("Unable to start charging connector")
				return
			}

			logInfo.Infof("Started charging connector at %s", time.Now())
		case types.AuthorizationStatusBlocked, types.AuthorizationStatusInvalid, types.AuthorizationStatusExpired:
			fallthrough
		default:
			logInfo.Errorf("Transaction unauthorized")
		}
	}

	err := util.SendRequest(cp.chargePoint, request, callback)
	if err != nil {

	}
}

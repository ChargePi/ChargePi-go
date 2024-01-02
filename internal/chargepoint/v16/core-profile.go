package v16

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/evse"
	"github.com/xBlaz3kx/ocppManager-go/ocpp_v16"
)

func (cp *ChargePoint) OnChangeAvailability(request *core.ChangeAvailabilityRequest) (confirmation *core.ChangeAvailabilityConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	response := core.AvailabilityStatusRejected

	// This would mean a request to change the availability of the whole charge point
	if request.ConnectorId == 0 {

		// Try to set the availability of the charge point, if it fails, schedule the change
		err := cp.SetAvailability(request.Type)
		if err != nil {
			_, err := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.SetAvailability, request.Type)
			if err != nil {
				return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusScheduled), nil
			}
		}

		return core.NewChangeAvailabilityConfirmation(core.AvailabilityStatusAccepted), nil
	}

	// Check if there are ongoing transactions, schedule the change if there are
	_, sessionErr := cp.sessionManager.GetSession(request.ConnectorId, nil)
	switch sessionErr {
	case nil:
		response = core.AvailabilityStatusScheduled
		// todo set evse availability
		// cp.evseManager.Get
	default:
		cp.logger.WithError(sessionErr).Error("Error checking for ongoing transactions ")
	}

	return core.NewChangeAvailabilityConfirmation(response), nil
}

func (cp *ChargePoint) OnChangeConfiguration(request *core.ChangeConfigurationRequest) (confirmation *core.ChangeConfigurationConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	var response = core.ConfigurationStatusRejected

	// Process the change configuration request
	switch request.Key {
	case ocpp_v16.AuthorizeRemoteTxRequests.String():
		// Just update
	case ocpp_v16.AllowOfflineTxForUnknownId.String():
		// Just update
	case ocpp_v16.AuthorizationCacheEnabled.String():
		cp.tagManager.ToggleAuthCache(request.Value == "true")
	case ocpp_v16.LocalAuthListMaxLength.String():
		val, err := strconv.Atoi(request.Value)
		if err != nil {
			return core.NewChangeConfigurationConfirmation(response), errors.New("invalid value")
		}

		cp.tagManager.SetMaxTags(val)
	case ocpp_v16.LocalAuthListEnabled.String():
		cp.tagManager.ToggleLocalAuthList(request.Value == "true")
	case ocpp_v16.SendLocalListMaxLength.String():
	case ocpp_v16.LightIntensity.String():
		brightness, err := strconv.Atoi(request.Value)
		if err != nil {
			return core.NewChangeConfigurationConfirmation(response), errors.New("invalid value")
		}

		err = cp.indicator.SetBrightness(brightness)
		if err != nil {
			return core.NewChangeConfigurationConfirmation(response), nil
		}

	case ocpp_v16.HeartbeatInterval.String():
		interval, err := strconv.Atoi(request.Value)
		if err != nil {
			return core.NewChangeConfigurationConfirmation(response), errors.New("invalid value")
		}

		cp.setHeartbeat(interval)
	case ocpp_v16.MeterValueSampleInterval.String():

	case ocpp_v16.MeterValuesSampledData.String():

	case ocpp_v16.AuthorizeRemoteTxRequests.String():
		// Just update
	default:
		return core.NewChangeConfigurationConfirmation(response), errors.New("unknown key")
	}

	err = cp.settingsManager.GetOcppV16Manager().UpdateKey(ocpp_v16.Key(request.Key), &request.Value)
	if err == nil {
		response = core.ConfigurationStatusAccepted
	}

	return core.NewChangeConfigurationConfirmation(response), nil
}

func (cp *ChargePoint) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	response := core.ClearCacheStatusRejected

	cacheErr := cp.tagManager.ClearCache()
	switch {
	case cacheErr == nil:
		cp.logger.Info("Cache cleared")
		response = core.ClearCacheStatusAccepted
	case errors.Is(cacheErr, auth.ErrCacheNotEnabled):
		cp.logger.Info("Cache not enabled")
		response = core.ClearCacheStatusRejected
	default:
		cp.logger.WithError(cacheErr).Warn("Unable to clear cache")
		response = core.ClearCacheStatusRejected
	}

	return core.NewClearCacheConfirmation(response), nil
}

func (cp *ChargePoint) OnDataTransfer(request *core.DataTransferRequest) (confirmation *core.DataTransferConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	var response = core.DataTransferStatusRejected

	// Supporting direct display control over custom data transfer messages, based on the messages in OCPP 2.0.1.
	if request.VendorId != cp.settingsManager.GetChargePointSettings().Info.OCPPDetails.Vendor {
		return core.NewDataTransferConfirmation(core.DataTransferStatusUnknownVendorId), nil
	}

	switch request.MessageId {
	case display.ClearDisplayMessageFeatureName:
		_ = request.Data.(display.ClearDisplayRequest)
		cp.display.Clear()
	case display.NotifyDisplayMessagesFeatureName:
		_ = request.Data.(display.NotifyDisplayMessagesRequest)
	case display.GetDisplayMessagesFeatureName:
		_ = request.Data.(display.GetDisplayMessagesRequest)
	case display.SetDisplayMessageFeatureName:
		req := request.Data.(display.SetDisplayMessageRequest)

		displayErr := cp.DisplayMessage(req.Message)
		if displayErr != nil {
			cp.logger.WithError(displayErr).Warn("Failed to display requested message")
			response = core.DataTransferStatusRejected
		}
	default:
		response = core.DataTransferStatusUnknownMessageId
	}

	return core.NewDataTransferConfirmation(response), nil
}

func (cp *ChargePoint) OnGetConfiguration(request *core.GetConfigurationRequest) (confirmation *core.GetConfigurationConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())

	var (
		unknownKeys            []string
		configArray            = []core.ConfigurationKey{}
		response               = core.NewGetConfigurationConfirmation(configArray)
		configuration, confErr = cp.settingsManager.GetOcppV16Manager().GetConfiguration()
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
		_, keyErr := cp.settingsManager.GetOcppV16Manager().GetConfigurationValue(ocpp_v16.Key(key))
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

	resetRetries, _ := cp.settingsManager.GetOcppV16Manager().GetConfigurationValue(ocpp_v16.ResetRetries)
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
	response := core.UnlockStatusNotSupported

	conn, fErr := cp.evseManager.GetEVSE(request.ConnectorId)
	switch fErr {
	case nil:
		cp.logger.Infof("Unlocking connector %d", request.ConnectorId)
		response = core.UnlockStatusUnlocked
		conn.GetEvcc().Unlock()
	}

	return core.NewUnlockConnectorConfirmation(response), nil
}

func (cp *ChargePoint) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	cp.logger.WithField("transactionId", request.TransactionId).Infof("Received request %s", request.GetFeatureName())

	response := types.RemoteStartStopStatusRejected
	transactionId := fmt.Sprintf("%d", request.TransactionId)

	session, fErr := cp.sessionManager.GetSessionWithTransactionId(transactionId)
	if fErr == nil {
		cp.logger.WithField("transactionId", request.TransactionId).Infof("Stopping transaction")
		response = types.RemoteStartStopStatusAccepted

		// Delay stopping the transaction by 3 seconds
		_, schedulerErr := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.StopCharging, session.EvseId, session.ConnectorId, core.ReasonRemote)
		if schedulerErr != nil {
			cp.logger.WithError(err).Error("Failed to schedule stop charging")
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
		conn, err = cp.evseManager.GetEVSE(*request.ConnectorId)
	} else {
		conn, err = cp.evseManager.GetAvailableEVSE()
	}

	if err == nil && conn.IsAvailable() {
		logInfo.Infof("Remote starting transaction")

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
	logInfo := cp.logger.WithFields(log.Fields{
		"evseId":      evseId,
		"connectorId": connectorId,
		"tagId":       tagId,
	})

	if cp.availability != core.AvailabilityTypeOperative {
		return
	}

	authorizeRemoteTx, _ := cp.settingsManager.GetOcppV16Manager().GetConfigurationValue(ocpp_v16.AuthorizeRemoteTxRequests)
	if authorizeRemoteTx != nil && *authorizeRemoteTx == "true" {
		logInfo.Info("Authorizing RemoteStart transaction")

		if !cp.isTagAuthorized(tagId) {
			logInfo.Warn("Tag unauthorized")
			return
		}
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
			measurements, sampleInterval := cp.getSessionParameters()
			err := cp.evseManager.StartCharging(evseId, nil, measurements, sampleInterval)
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

	err := cp.sendRequest(request, callback)
	if err != nil {

	}
}

package v16

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/samber/lo"
)

func (cp *ChargePoint) StartCharging(evseId, connectorId int, tagId string) error {
	logger := cp.logger.With(zap.Int("evseId", evseId), zap.Int("connectorId", connectorId), zap.String("tagId", tagId))
	logger.Info("Starting charging")

	// Charge point must be available to accept transactions
	if cp.availability != core.AvailabilityTypeOperative {
		return chargePoint.ErrChargePointUnavailable
	}

	// Authorize the tag from either the list, cache or the backend.
	if !cp.isTagAuthorized(tagId) {
		return chargePoint.ErrTagUnauthorized
	}

	// todo sample power meter to get current energy value

	request := core.NewStartTransactionRequest(
		evseId,
		tagId,
		0,
		types.NewDateTime(time.Now()),
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logger.With(zap.Error(protoError)).Sugar().Warnf("Central system responded with an error for %s", confirmation.GetFeatureName())
			return
		}

		startTransactionConf := confirmation.(*core.StartTransactionConfirmation)

		switch startTransactionConf.IdTagInfo.Status {
		case types.AuthorizationStatusAccepted, types.AuthorizationStatusConcurrentTx:
			transactionId := fmt.Sprintf("%d", startTransactionConf.TransactionId)

			err := cp.sessionManager.StartSession(evseId, nil, tagId, transactionId)
			if err != nil {
				logger.With(zap.Error(err)).Error("Unable to start a session")
				return
			}

			// Get metering parameters
			measurands, sampleInterval := cp.getSessionParameters()

			// Start the charging on EVSE
			err = cp.evseManager.StartCharging(evseId, nil, measurands, sampleInterval)
			if err != nil {
				logger.With(zap.Error(err)).Error("Unable to start charging on EVSE")
				return
			}

			logger.Sugar().Infof("Started charging connector at %s", time.Now())
		case types.AuthorizationStatusBlocked, types.AuthorizationStatusInvalid, types.AuthorizationStatusExpired:
			fallthrough
		default:
			logger.Warn("Transaction unauthorized")
		}

		// Cache the tag
		err := cp.tagManager.AddTag(tagId, startTransactionConf.IdTagInfo)
		if err != nil {
			logger.With(zap.Error(err)).Warn("Unable to cache tag")
		}
	}

	return cp.sendRequest(request, callback)
}

func (cp *ChargePoint) getSessionParameters() ([]types.Measurand, string) {
	cp.logger.Debug("Getting session sampling parameters")

	// Get metering parameters
	variableManager := cp.settingsManager.GetOcppV16Manager()
	sampleInterval, err := variableManager.GetConfigurationValue(ocpp_v16.MeterValueSampleInterval)
	if err != nil {
		sampleInterval = lo.ToPtr("90s")
	}

	measurandsString, err := variableManager.GetConfigurationValue(ocpp_v16.MeterValuesSampledData)
	if err != nil {
		measurandsString = lo.ToPtr(string(types.MeasurandEnergyActiveImportRegister))
	}

	var measurands []types.Measurand
	for _, measurand := range strings.Split(*measurandsString, ",") {
		measurands = append(measurands, types.Measurand(measurand))
	}

	return measurands, *sampleInterval
}

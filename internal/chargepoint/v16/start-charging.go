package v16

import (
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ocpp-manager/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (cp *ChargePoint) OnRemoteStartTransaction(request *core.RemoteStartTransactionRequest) (confirmation *core.RemoteStartTransactionConfirmation, err error) {
	var (
		logger = cp.logger.With(
			zap.Intp("connectorId", request.ConnectorId),
			zap.String("tagId", request.IdTag),
		)
		response = types.RemoteStartStopStatusRejected
		conn     evse.EVSE
	)

	logger.Sugar().Infof("Received request %s", request.GetFeatureName())

	// If the connector is specified, check if it exists and is available.
	if request.ConnectorId != nil {
		conn, err = cp.evseManager.GetEVSE(*request.ConnectorId)
	} else {
		conn, err = cp.evseManager.GetAvailableEVSE()
	}

	if err == nil && conn.IsAvailable() {
		logger.Info("Remote starting transaction")

		// Delay the charging by 3 seconds
		response = types.RemoteStartStopStatusAccepted
		_, schedulerErr := cp.scheduler.Every(1).Seconds().LimitRunsTo(1).Do(cp.remoteStart, conn, 1, request.IdTag)
		if schedulerErr != nil {
			response = types.RemoteStartStopStatusRejected
		}
	}

	return core.NewRemoteStartTransactionConfirmation(response), nil
}

func (cp *ChargePoint) remoteStart(evseId, connectorId int, tagId string) {
	logger := cp.logger.With(
		zap.Int("evseId", evseId),
		zap.Int("connectorId", connectorId),
		zap.String("tagId", tagId),
	)

	if cp.state.GetAvailability() != core.AvailabilityTypeOperative {
		return
	}

	// Authorize the tag from either the list, cache or the backend.
	authorizeRemoteTx, _ := cp.settingsManager.GetConfigurationValue(ocpp_v16.AuthorizeRemoteTxRequests)
	if authorizeRemoteTx != nil && *authorizeRemoteTx == "true" {
		logger.Info("Authorizing RemoteStart transaction")

		if !cp.isAuthorized(tagId) {
			logger.Warn("Tag unauthorized")
			return
		}
	}

	// Start the charging
	err := cp.StartCharging(evseId, connectorId, tagId)
	if err != nil {
		logger.With(zap.Error(err)).Error("Unable to start charging remotely")
	}
}

func (cp *ChargePoint) StartCharging(evseId, connectorId int, tagId string) error {
	logger := cp.logger.With(zap.Int("evseId", evseId), zap.Int("connectorId", connectorId), zap.String("tagId", tagId))
	logger.Info("Starting charging")

	// Charge point must be available to accept transactions
	if cp.state.GetAvailability() != core.AvailabilityTypeOperative {
		return chargepoint.ErrChargePointUnavailable
	}

	// Authorize the tag from either the list, cache or the backend. If no authorization is required, skip this step.
	if !cp.isAuthorized(tagId) {
		return chargepoint.ErrTagUnauthorized
	}

	// Get sampling parameters
	measurands, sampleInterval := cp.getSamplingParameters()

	// Start the charging on EVSE
	err := cp.evseManager.StartCharging(evseId, nil, measurands, sampleInterval)
	if err != nil {
		logger.With(zap.Error(err)).Error("Unable to start charging on EVSE")
		return err
	}

	// Log the session in the session manager
	err = cp.sessionService.StartSession(evseId, nil, tagId, "")
	if err != nil {
		logger.With(zap.Error(err)).Error("Unable to start a session")
		return err
	}

	logger.Sugar().Infof("Started charging connector at %s", time.Now())

	// Notify the central system that the transaction has started
	request := core.NewStartTransactionRequest(
		evseId,
		tagId,
		0, // todo sample power meter to get current energy value
		types.NewDateTime(time.Now()),
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logger.With(zap.Error(protoError)).Sugar().Warnf("Central system responded with an error for %s", confirmation.GetFeatureName())
			return
		}

		startTransactionConf := confirmation.(*core.StartTransactionConfirmation)

		// Update the transaction id in the session
		logger.Sugar().Infof("Updating transaction ID: %d", startTransactionConf.TransactionId)
		err = cp.sessionService.AddTransactionIdToSession(evseId, nil, strconv.Itoa(startTransactionConf.TransactionId))
		if err != nil {
			logger.With(zap.Error(err)).Warn("Unable to update transaction ID in the session manager")
		}

		// Cache the tag
		err = cp.tagAuthService.CacheTag(tagId, startTransactionConf.IdTagInfo)
		if err != nil {
			logger.With(zap.Error(err)).Warn("Unable to cache tag")
		}

		defer cp.isAuthorized(tagId)
	}

	return cp.sendRequest(request, callback)
}

// StartChargingFreeMode starts a charging session in free mode.
func (cp *ChargePoint) StartChargingFreeMode(evseId int) error {
	if !cp.info.FreeMode {
		return errors.New("free mode is not enabled")
	}

	logInfo := cp.logger.With(zap.Int("evseId", evseId))
	logInfo.Info("Free mode enabled, starting charging")

	// Get sampling parameters
	measurements, sampleInterval := cp.getSamplingParameters()

	// Start the charging on EVSE
	return cp.evseManager.StartCharging(evseId, nil, measurements, sampleInterval)
}

// getSamplingParameters retrieves the sampling parameters from the charge point settings.
func (cp *ChargePoint) getSamplingParameters() ([]types.Measurand, string) {
	cp.logger.Debug("Getting session sampling parameters")

	// Get metering parameters
	sampleInterval := cp.getSamplingInterval()
	measurands := cp.getMeasurands()

	return measurands, *sampleInterval
}

func (cp *ChargePoint) getSamplingInterval() *string {
	sampleInterval, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.MeterValueSampleInterval)
	if err != nil {
		// Default to 90 seconds
		sampleInterval = lo.ToPtr("90s")
	}

	return sampleInterval
}

func (cp *ChargePoint) getMeasurands() []types.Measurand {
	// Get measurands to sample
	measurandsString, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.MeterValuesSampledData)
	if err != nil {
		measurandsString = lo.ToPtr(string(types.MeasurandEnergyActiveImportRegister))
	}

	var measurands []types.Measurand
	for _, measurand := range strings.Split(*measurandsString, ",") {
		measurands = append(measurands, types.Measurand(measurand))
	}

	return measurands
}

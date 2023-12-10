package evse

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/power-meter"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

func (evse *Impl) GetPowerMeter() powerMeter.PowerMeter {
	return evse.powerMeter
}

func (evse *Impl) SetPowerMeter(meter powerMeter.PowerMeter) error {
	evse.powerMeter = meter
	return nil
}

// SamplePowerMeter Get a sample from the power meter. The measurands argument takes the list of all the types of the measurands to sample.
// It will add all the samples to the evse's Session if it is active.
func (evse *Impl) SamplePowerMeter(measurands []types.Measurand) []types.SampledValue {
	logInfo := evse.logger

	if util.IsNilInterfaceOrPointer(evse.powerMeter) {
		logInfo.Warn("Sampling the power meter unavailable")
		return nil
	}

	logInfo.Debugf("Sampling EVSE for measurands %v", measurands)

	var samples []types.SampledValue

	// Get value for each supported measureand
	for _, measurand := range measurands {
		logInfo.Debugf("Sampling measurand %v", measurand)

		switch measurand {
		case types.MeasurandPowerActiveImport, types.MeasurandPowerActiveExport:
			power, err := evse.powerMeter.GetPower(1)
			if err != nil {
				continue
			}
			samples = append(samples, *power)
		case types.MeasurandEnergyActiveImportInterval, types.MeasurandEnergyActiveImportRegister,
			types.MeasurandEnergyActiveExportInterval, types.MeasurandEnergyActiveExportRegister:
			energy, err := evse.powerMeter.GetEnergy()
			if err != nil {
				continue
			}

			samples = append(samples, *energy)
		case types.MeasurandCurrentImport, types.MeasurandCurrentExport:
			current, err := evse.powerMeter.GetCurrent(1)
			if err != nil {
				continue
			}
			samples = append(samples, *current)

			current, err = evse.powerMeter.GetCurrent(2)
			if err != nil {
				continue
			}
			samples = append(samples, *current)

			current, err = evse.powerMeter.GetCurrent(3)
			if err != nil {
				continue
			}
			samples = append(samples, *current)
		case types.MeasurandVoltage:
			voltage, err := evse.powerMeter.GetVoltage(1)
			if err != nil {
				continue
			}
			samples = append(samples, *voltage)

			voltage, err = evse.powerMeter.GetVoltage(2)
			if err != nil {
				continue
			}
			samples = append(samples, *voltage)

			voltage, err = evse.powerMeter.GetVoltage(3)
			if err != nil {
				continue
			}
			samples = append(samples, *voltage)
		}
	}

	return samples
}

func (evse *Impl) samplePowerMeterAndSend() {
	evse.logger.Debug("Sampling power meter")

	measurands := util.GetTypesToSample()
	samples := evse.SamplePowerMeter(measurands)

	meterValue := types.MeterValue{
		Timestamp:    types.NewDateTime(time.Now()),
		SampledValue: samples,
	}

	// Notify a MeterValue update
	if evse.meterValuesChannel != nil {
		evse.logger.Debugf("Sending meter value notification")
		// todo get transaction id evse.GetTransactionId()
		evse.meterValuesChannel <- notifications.NewMeterValueNotification(evse.evseId, &evse.evseId, nil, meterValue)
	}
}

// preparePowerMeter
func (evse *Impl) preparePowerMeter() error {
	if util.IsNilInterfaceOrPointer(evse.powerMeter) {
		return ErrPowerMeterNotEnabled
	}

	evse.logger.Debug("Preparing power meter")

	var (
		sampleTime          = "30s"
		sampleInterval, err = ocppConfigManager.GetConfigurationValue(configuration.MeterValueSampleInterval.String())
	)

	if err == nil && sampleInterval != nil {
		sampleTime = fmt.Sprintf("%ss", *sampleInterval)
	}

	// Schedule the sampling
	_, err = evse.scheduler.Every(sampleTime).
		Tag("evse", "sampling", fmt.Sprintf("%d", evse.evseId)).
		Do(evse.samplePowerMeterAndSend)

	return err
}

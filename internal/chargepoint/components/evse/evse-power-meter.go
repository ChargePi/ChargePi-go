package evse

import (
	"fmt"
	"time"

	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/pkg/power-meter"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
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
func (evse *Impl) SamplePowerMeter(measurands []types.Measurand) {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
	})

	if !evse.powerMeterEnabled || util.IsNilInterfaceOrPointer(evse.powerMeter) {
		logInfo.Warn("Sampling the power meter unavailable")
		return
	}

	logInfo.Debugf("Sampling EVSE with: %v", measurands)

	var (
		meterValues = types.MeterValue{SampledValue: []types.SampledValue{}, Timestamp: types.NewDateTime(time.Now())}
		samples     []types.SampledValue
	)

	// Get value for each supported measureand
	for _, measurand := range measurands {
		switch measurand {
		case types.MeasurandPowerActiveImport, types.MeasurandPowerActiveExport:
			samples = append(samples, evse.powerMeter.GetPower())
		case types.MeasurandEnergyActiveImportInterval, types.MeasurandEnergyActiveImportRegister,
			types.MeasurandEnergyActiveExportInterval, types.MeasurandEnergyActiveExportRegister:
			samples = append(samples, evse.powerMeter.GetEnergy())
		case types.MeasurandCurrentImport, types.MeasurandCurrentExport:
			samples = append(samples, evse.powerMeter.GetCurrent(1))
			samples = append(samples, evse.powerMeter.GetCurrent(2))
			samples = append(samples, evse.powerMeter.GetCurrent(3))
		case types.MeasurandVoltage:
			samples = append(samples, evse.powerMeter.GetVoltage(1))
			samples = append(samples, evse.powerMeter.GetVoltage(2))
			samples = append(samples, evse.powerMeter.GetVoltage(3))
		}
	}

	meterValues.SampledValue = samples

	// Notify a MeterValue update
	if evse.meterValuesChannel != nil {
		// evse.GetTransactionId()
		evse.meterValuesChannel <- notifications.NewMeterValueNotification(evse.evseId, &evse.evseId, nil, meterValues)
	}

	evse.session.AddSampledValue(samples)
}

// preparePowerMeter
func (evse *Impl) preparePowerMeter() error {
	if !evse.powerMeterEnabled || util.IsNilInterfaceOrPointer(evse.powerMeter) {
		return ErrPowerMeterNotEnabled
	}

	var (
		measurands          = util.GetTypesToSample()
		sampleTime          = "30s"
		sampleInterval, err = ocppConfigManager.GetConfigurationValue(configuration.MeterValueSampleInterval.String())
		jobTag              = fmt.Sprintf("Evse%dSampling", evse.evseId)
	)

	if err == nil && sampleInterval != nil {
		sampleTime = fmt.Sprintf("%ss", *sampleInterval)
	}

	// Schedule the sampling
	_, err = scheduler.NewScheduler().Every(sampleTime).
		Tag(jobTag).
		Do(evse.SamplePowerMeter, measurands)

	return err
}

func (evse *Impl) CalculateSessionAvgEnergyConsumption() float64 {
	return evse.session.CalculateEnergyConsumptionWithAvgPower()
}

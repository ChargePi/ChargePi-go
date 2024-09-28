//go:build linux || raspberrypi4

package evse

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/ChargePi/ChargePi-go/internal/pkg/scheduler"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/evcc"
	powerMeter "github.com/ChargePi/ChargePi-go/pkg/hardware/power-meter"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	EvseId     int                 `json:"evseId,omitempty" yaml:"evseId" mapstructure:"evseId" validate:"required,gte=1"`
	MaxPower   float32             `json:"maxPower" yaml:"maxPower" mapstructure:"maxPower" validate:"gt=0"`
	EVCC       evcc.Settings       `json:"evcc" yaml:"evcc" mapstructure:"evcc" validate:"required"`
	PowerMeter powerMeter.Settings `json:"powerMeter" yaml:"powerMeter" mapstructure:"powerMeter"`
	Connectors []ConnectorSettings `json:"connectors" yaml:"connectors" mapstructure:"connectors"`
}

func (s *Settings) Validate() error {
	return validator.New().Struct(s)
}

type V1 struct {
	// Settings
	info       Info
	connectors []ConnectorSettings

	scheduler *gocron.Scheduler
	logger    log.FieldLogger

	// State
	state *state

	// Notification channels
	meterValuesChannel  chan<- notifications.MeterValueNotification
	notificationChannel chan<- notifications.StatusNotification

	// Hardware
	powerMeterEnabled bool
	powerMeter        powerMeter.PowerMeter
	evcc              evcc.EVCC
}

// NewEvse Creates a new evse.
func NewEvse(evseId int, evcc evcc.EVCC, powerMeter powerMeter.PowerMeter, maxPower float64) (*V1, error) {
	log.WithFields(log.Fields{
		"evseId": evseId,
	}).Info("Creating a new evse")

	if evseId <= 0 {
		return nil, ErrInvalidEvseId
	}

	if util.IsNilInterfaceOrPointer(evcc) {
		return nil, ErrInvalidEVCC
	}

	powerMeterEnabled := false
	if util.IsNilInterfaceOrPointer(powerMeter) {
		powerMeterEnabled = true
	}

	return &V1{
		info: Info{
			EvseId:   evseId,
			MaxPower: maxPower,
		},
		evcc:              evcc,
		powerMeter:        powerMeter,
		powerMeterEnabled: powerMeterEnabled,
		connectors:        make([]ConnectorSettings, 0),
		state:             newState(),
		scheduler:         scheduler.NewScheduler(),
		logger:            log.StandardLogger().WithField("component", "evse").WithField("evseId", evseId),
	}, nil
}

// NewEvseFromSettings Create a new evse object from the configuration.
func NewEvseFromSettings(settings Settings) (*V1, error) {
	logInfo := log.StandardLogger()

	err := settings.Validate()
	if err != nil {
		return nil, err
	}

	// Create an EVCC from settings
	evccFromType, err := evcc.NewEVCCFromType(settings.EVCC)
	switch err {
	case nil:
		logInfo.WithField("type", settings.EVCC.Type).Debugf("EVCC created")
	default:
		return nil, err
	}

	// Create a PowerMeter from settings
	logInfo.Debugf("Creating power meter")
	meter, powerMeterErr := powerMeter.NewPowerMeter(settings.PowerMeter)
	switch {
	case powerMeterErr == nil:
	case errors.Is(powerMeterErr, powerMeter.ErrPowerMeterDisabled):
		logInfo.WithError(powerMeterErr).Warn("Power meter disabled")
	case errors.Is(powerMeterErr, powerMeter.ErrPowerMeterUnsupported), errors.Is(powerMeterErr, powerMeter.ErrInvalidConnectionSettings):
		fallthrough
	default:
		logInfo.WithError(powerMeterErr).Error("Cannot instantiate power meter for evse")
		return nil, powerMeterErr
	}

	return &V1{
		info: Info{
			EvseId:   settings.EvseId,
			MaxPower: float64(settings.MaxPower),
		},
		evcc:       evccFromType,
		powerMeter: meter,
		state:      newState(),
		scheduler:  scheduler.NewScheduler(),
		logger:     log.StandardLogger().WithField("component", "evse").WithField("evseId", settings.EvseId),
	}, nil
}

// Init Initialize the evse.This method must be called before using the EVSE.
func (evse *V1) Init(ctx context.Context) error {
	evse.logger.Info("Initializing evse")
	// Init EVCC
	err := evse.evcc.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot initialize evcc")
	}

	// Disable charging by default
	evse.evcc.DisableCharging()

	// Set max charging current
	err = evse.evcc.SetMaxChargingCurrent(evse.info.MaxPower)
	if err != nil {
		return err
	}

	// Listen for EVCC status updates in another thread.
	go func() {
		defer evse.Cleanup()
		evse.listenForStatusUpdates(ctx)
	}()

	return nil
}

func (evse *V1) Cleanup() error {
	evse.logger.Info("Cleaning up evse")

	// Stop charging
	err := evse.StopCharging(core.ReasonLocal)
	if err != nil {
		evse.logger.WithError(err).Error("Cannot stop charging")
	}

	// Clean up EVCC
	err = evse.evcc.Cleanup()
	if err != nil {
		evse.logger.WithError(err).Error("Cannot cleanup evcc")
	}

	// Stop the scheduler
	evse.scheduler.Stop()

	// Clean up the power meter
	if evse.powerMeterEnabled {
		err = evse.powerMeter.Cleanup()
		if err != nil {
			evse.logger.WithError(err).Error("Cannot stop power meter")
		}
	}

	return nil
}

// listenForStatusUpdates Listen for status updates from the evcc and update the evse status accordingly.
func (evse *V1) listenForStatusUpdates(ctx context.Context) {
	evse.logger.Debug("Listening for evcc status updates")

	statusChan := evse.evcc.GetStatusChangeChannel()
	if statusChan == nil {
		evse.logger.Panic("Cannot listen for evcc status updates")
	}

Loop:
	for {
		select {
		case msg := <-statusChan:
			// Determine OCPP status based on CarState and Error

			var (
				state core.ChargePointStatus
				cpErr core.ChargePointErrorCode
			)

			currentStatus, _ := evse.GetStatus()

			// Compare to current status
			switch currentStatus {
			case core.ChargePointStatusAvailable:

				switch msg.State {
				case evcc.StateB1:
					state = core.ChargePointStatusPreparing
				}

			case core.ChargePointStatusPreparing:

				// Determine new state based on the previous state
				switch msg.State {
				case evcc.StateA1, evcc.StateA2:
					state = core.ChargePointStatusAvailable
				case evcc.StateC2, evcc.StateD2:
					state = core.ChargePointStatusCharging
				}

			case core.ChargePointStatusCharging:

				switch msg.State {
				case evcc.StateC1:
					state = core.ChargePointStatusFinishing
				case evcc.StateD1:
					state = core.ChargePointStatusSuspendedEV
				}

				// todo
			case core.ChargePointStatusSuspendedEV:

				switch msg.State {
				case evcc.StateC1:
					state = core.ChargePointStatusFinishing
				case evcc.StateB1:
					state = core.ChargePointStatusPreparing
				}

			case core.ChargePointStatusSuspendedEVSE:

				switch msg.State {
				case evcc.StateC2:
					state = core.ChargePointStatusCharging
				case evcc.StateB1:
					state = core.ChargePointStatusFinishing
				case evcc.StateA1:
					state = core.ChargePointStatusAvailable
				}

			case core.ChargePointStatusFaulted:
			}

			switch msg.State {
			case evcc.StateE, evcc.StateF:
				state = core.ChargePointStatusFaulted
			}

			err := evse.SetStatus(state, cpErr)
			if err != nil {
				// todo is this a panic event?
				evse.logger.WithError(err).Error("Cannot set evse status")
			}
		case <-ctx.Done():
			break Loop
		}
	}
}

func (evse *V1) SetState(state State) {

}

// StartCharging Start charging an evse if evse is available and session could be started.
func (evse *V1) StartCharging(connectorId *int, measurands []types.Measurand, sampleInterval string) error {
	logInfo := evse.logger.WithField("connectorId", connectorId)
	logInfo.Debugf("Trying to start charging on evse")

	// Enable charging on the evcc
	err := evse.evcc.EnableCharging()
	if err != nil {
		return err
	}

	// Lock the connector
	evse.evcc.Lock()

	// Prepare power meter and schedule sampling
	sampleError := evse.scheduleMeterValueUpdates(measurands, sampleInterval)
	if sampleError != nil {
		logInfo.WithError(sampleError).Error("Cannot sample evse")
	}

	return nil
}

// StopCharging Stops charging an evse if evse is charging
func (evse *V1) StopCharging(reason core.Reason) error {
	logInfo := evse.logger.WithField("reason", reason)

	if evse.IsCharging() {
		logInfo.Debugf("Stopping charging")

		evse.evcc.DisableCharging()
		evse.evcc.Unlock()

		// Remove any jobs scheduled for this evse
		schedulerErr := evse.scheduler.RemoveByTag(fmt.Sprintf("evse-%d-chargingTimer", evse.GetEvseId()))
		if schedulerErr != nil {
			logInfo.WithError(schedulerErr).Errorf("Cannot remove sampling schedule")
		}

		return nil
	}

	return ErrNotCharging
}

func (evse *V1) IsAvailable() bool {
	availability := evse.state.GetAvailability()
	status, _ := evse.state.GetStatus()

	return status == core.ChargePointStatusAvailable && availability == core.AvailabilityTypeOperative
}

func (evse *V1) IsCharging() bool {
	status, _ := evse.state.GetStatus()
	return status == core.ChargePointStatusCharging
}

func (evse *V1) IsPreparing() bool {
	status, _ := evse.state.GetStatus()
	return status == core.ChargePointStatusPreparing
}

func (evse *V1) IsReserved() bool {
	status, _ := evse.state.GetStatus()
	return status == core.ChargePointStatusReserved
}

func (evse *V1) IsUnavailable() bool {
	status, _ := evse.state.GetStatus()
	return status == core.ChargePointStatusUnavailable
}

// SetAvailability sets the availability of the evse.
func (evse *V1) SetAvailability(isAvailable bool) error {
	// Check if state has changed
	if evse.state.GetAvailability() == core.AvailabilityTypeOperative && isAvailable {
		return nil
	}

	if evse.state.GetAvailability() == core.AvailabilityTypeInoperative && !isAvailable {
		return nil
	}

	if isAvailable {
		return evse.state.SetAvailability(core.AvailabilityTypeOperative)
	}

	err := evse.state.SetAvailability(core.AvailabilityTypeInoperative)
	if err != nil {
		return err
	}

	// Disable charging immediately if the evse is not available
	evse.evcc.DisableCharging()
	return nil
}

// SetStatus sets the status of the EVSE and sends a notification to the channel
func (evse *V1) SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) error {
	logInfo := evse.logger.WithFields(log.Fields{
		"status": status,
		"err":    errCode,
	})
	logInfo.Debugf("Setting evse status %s with err %s", status, errCode)

	err := evse.state.SetStatus(status, errCode)
	if err != nil {
		return err
	}

	// Notify the channel that a status was updated
	if evse.notificationChannel != nil {
		logInfo.Debug("Sending status notification")
		evse.notificationChannel <- notifications.NewStatusNotification(evse.info.EvseId, string(status), string(errCode))
	}
	return nil
}

// GetStatus returns the status of the EVSE
func (evse *V1) GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode) {
	return evse.state.GetStatus()
}

// GetEvseId returns the id of the evse
func (evse *V1) GetEvseId() int {
	evse.logger.Debugf("Getting evse id")
	return evse.info.EvseId
}

// GetMaxChargingPower returns the maximum charging power of the evse
func (evse *V1) GetMaxChargingPower() float64 {
	evse.logger.Debugf("Getting max charging power")
	return evse.info.MaxPower
}

// SetNotificationChannel sets the notification channel for the evse
func (evse *V1) SetNotificationChannel(notificationChannel chan<- notifications.StatusNotification) {
	evse.notificationChannel = notificationChannel
}

// IsHealthy checks if underlying components are behaving as expected
func (evse *V1) IsHealthy() bool {
	if evse.powerMeterEnabled && !util.IsNilInterfaceOrPointer(evse.powerMeter) {
		// Voltage should always be available
		_, err := evse.powerMeter.GetVoltage(1)
		if err != nil {
			return false
		}
	}

	// Check if EVCC healthy
	return evse.evcc.GetError() == ""
}

// Lock locks the connector with the given id
func (evse *V1) Lock(connectorId int) {
	evse.logger.Debug("Locking EVCC")
	evse.evcc.Lock()
}

// Unlock unlocks the connector with the given id
func (evse *V1) Unlock(connectorId int) {
	evse.logger.Debug("Unlocking EVCC")
	evse.evcc.Unlock()
}

// GetConnectors returns all connectors ascending by id.
func (evse *V1) GetConnectors() []ConnectorSettings {
	evse.logger.Debug("Getting connectors for EVSE")

	return evse.connectors
}

// AddConnector adds a connector to the EVSE. The connector id must be unique.
// It will reorder the connectors based on the ID.
func (evse *V1) AddConnector(connector ConnectorSettings) error {
	// Validate the settings
	err := connector.Validate()
	if err != nil {
		return err
	}

	evse.logger.WithField("connectorId", connector.ConnectorId).Debug("Adding connector to EVSE")

	containsConnectorWithId := lo.ContainsBy(evse.connectors, func(item ConnectorSettings) bool {
		return item.ConnectorId == connector.ConnectorId
	})
	if containsConnectorWithId {
		return ErrConnectorExists
	}

	// Append to the array and sort by connector id
	evse.connectors = append(evse.connectors, connector)
	sort.Slice(evse.connectors, func(i, j int) bool {
		return evse.connectors[i].ConnectorId < evse.connectors[j].ConnectorId
	})
	return nil
}

func (evse *V1) GetEvcc() evcc.EVCC {
	evse.logger.Debug("Getting EVCC")
	return evse.evcc
}

// SetEvcc modifies the EVCC component at runtime,
// cleaning up the previous running EVCC and initializing the new one.
func (evse *V1) SetEvcc(e evcc.EVCC) error {
	evse.logger.Debug("Setting EVCC")

	if util.IsNilInterfaceOrPointer(e) {
		return ErrInvalidEVCC
	}

	// Cleanup the "current" EVCC
	err := evse.evcc.Cleanup()
	if err != nil {
		evse.logger.Errorf("Error cleaning up EVCC: %s", err)
		return err
	}

	// Run the init function
	// e.Init(nil)

	evse.evcc = e
	return nil
}

func (evse *V1) SetMeterValuesChannel(notificationChannel chan<- notifications.MeterValueNotification) {
	evse.meterValuesChannel = notificationChannel
}

func (evse *V1) GetPowerMeter() powerMeter.PowerMeter {
	return evse.powerMeter
}

// SetPowerMeter sets the power meter for the EVSE at runtime.
// It cleans up the previous power meter and initialize the new one.
func (evse *V1) SetPowerMeter(meter powerMeter.PowerMeter) error {
	if util.IsNilInterfaceOrPointer(meter) {
		return errors.New("power meter cannot be nil")
	}

	evse.logger.Debug("Setting power meter")
	if !util.IsNilInterfaceOrPointer(evse.powerMeter) {
		// Cleanup previous power meter
		err := evse.powerMeter.Cleanup()
		if err != nil {
			return err
		}
	}

	// Set new power meter
	evse.powerMeter = meter
	evse.powerMeterEnabled = true
	return nil
}

// SamplePowerMeter requests samples from the power meter for given measurands.
func (evse *V1) SamplePowerMeter(measurands []types.Measurand) ([]types.SampledValue, error) {
	logInfo := evse.logger

	if util.IsNilInterfaceOrPointer(evse.powerMeter) {
		logInfo.Warn("Sampling the power meter unavailable")
		return nil, errors.New("power meter not enabled")
	}

	logInfo.Debugf("Sampling EVSE for measurands %v", measurands)

	var samples []types.SampledValue

	// Get value for each supported measureand
	for _, measurand := range measurands {
		logInfo.Debugf("Sampling measurand %v", measurand)

		switch measurand {
		case types.MeasurandPowerActiveImport, types.MeasurandPowerActiveExport:
			// Get the total power
			current, err := evse.powerMeter.GetPower(1)
			if err != nil {
				logInfo.WithError(err).Error("Error sampling power meter power")
				continue
			}

			samples = append(samples, *current)
		case types.MeasurandEnergyActiveImportInterval, types.MeasurandEnergyActiveImportRegister,
			types.MeasurandEnergyActiveExportInterval, types.MeasurandEnergyActiveExportRegister:
			energy, err := evse.powerMeter.GetEnergy()
			if err != nil {
				logInfo.WithError(err).Error("Error sampling power meter energy")
				continue
			}

			samples = append(samples, *energy)
		case types.MeasurandCurrentImport, types.MeasurandCurrentExport:
			// Get current for each phase
			for i := 1; i < 4; i++ {
				current, err := evse.powerMeter.GetCurrent(i)
				if err != nil {
					logInfo.WithError(err).Error("Error sampling power meter current")
					continue
				}

				samples = append(samples, *current)
			}
		case types.MeasurandVoltage:
			// Get voltage for each phase
			for i := 1; i < 4; i++ {
				current, err := evse.powerMeter.GetVoltage(i)
				if err != nil {
					logInfo.WithError(err).Error("Error sampling power meter voltage")
					continue
				}

				samples = append(samples, *current)
			}
		}
	}

	return samples, nil
}

// sendMeterValueUpdate samples the power meter and sends a notification to the meterValuesChannel.
func (evse *V1) sendMeterValueUpdate(measurands []types.Measurand) {
	samples, err := evse.SamplePowerMeter(measurands)
	if err != nil {
		evse.logger.WithError(err).Error("Error sampling power meter")
		return
	}

	meterValue := types.MeterValue{
		Timestamp:    types.NewDateTime(time.Now()),
		SampledValue: samples,
	}

	// Notify a MeterValue update
	if evse.meterValuesChannel != nil {
		evse.logger.Debugf("Sending meter value notification")
		evse.meterValuesChannel <- notifications.NewMeterValueNotification(evse.info.EvseId, nil, nil, meterValue)
	}
}

// scheduleMeterValueUpdates schedules the sampling of the power meter at a given interval, if the power meter is enabled.
func (evse *V1) scheduleMeterValueUpdates(measurands []types.Measurand, samplingInterval string) error {
	if util.IsNilInterfaceOrPointer(evse.powerMeter) {
		return ErrPowerMeterNotEnabled
	}

	evse.logger.Debug("Preparing power meter")

	// Schedule the sampling
	_, err := evse.scheduler.Every(samplingInterval).
		Tag("evse", "sampling", fmt.Sprintf("%d", evse.info.EvseId)).
		Do(evse.sendMeterValueUpdate, measurands)

	return err
}

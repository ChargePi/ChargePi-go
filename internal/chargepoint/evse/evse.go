package evse

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-co-op/gocron"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/evcc"
	"github.com/xBlaz3kx/ChargePi-go/pkg/power-meter"
)

var (
	ErrInvalidEvseId        = errors.New("invalid evse id")
	ErrInvalidStatus        = errors.New("invalid evse status")
	ErrInvalidEVCC          = errors.New("evcc cannot be nil")
	ErrNotCharging          = errors.New("evse not charging")
	ErrPowerMeterNotEnabled = errors.New("power meter not enabled")
	ErrConnectorExists      = errors.New("connector already exists")
)

type (
	EVSE interface {
		Init(ctx context.Context) error

		AddConnector(connector settings.Connector) error
		GetConnectors() []settings.Connector
		GetMaxChargingTime() *int
		SetMaxChargingTime(time *int)
		GetMaxChargingPower() float64

		StartCharging(connectorId *int, measurands []types.Measurand, sampleInterval string) error
		StopCharging(reason core.Reason) error
		Lock()
		Unlock()

		GetEvseId() int

		GetPowerMeter() powerMeter.PowerMeter
		SetPowerMeter(powerMeter.PowerMeter) error
		SamplePowerMeter(measurands []types.Measurand) []types.SampledValue

		GetEvcc() evcc.EVCC
		SetEvcc(evcc.EVCC)

		SetAvailability(isAvailable bool)
		SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode)
		GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode)

		IsAvailable() bool
		IsPreparing() bool
		IsCharging() bool
		IsReserved() bool
		IsUnavailable() bool

		SetNotificationChannel(notificationChannel chan<- notifications.StatusNotification)
		SetMeterValuesChannel(notificationChannel chan<- notifications.MeterValueNotification)
	}

	Impl struct {
		evseId          int
		maxPower        float64
		maxChargingTime *int
		connectors      []settings.Connector
		availability    core.AvailabilityType
		status          core.ChargePointStatus
		errorCode       core.ChargePointErrorCode
		reservationId   *int

		// Notification channels
		meterValuesChannel  chan<- notifications.MeterValueNotification
		notificationChannel chan<- notifications.StatusNotification
		mu                  sync.Mutex
		scheduler           *gocron.Scheduler

		// Hardware
		powerMeterEnabled bool
		powerMeter        powerMeter.PowerMeter
		evcc              evcc.EVCC

		logger log.FieldLogger
	}
)

// NewEvse Create a new evse object from the provided arguments. evseId, connectorId and maxChargingTime must be greater than zero.
// When created, it makes an empty session, turns off the relay and defaults the status to Available.
func NewEvse(evseId int, evcc evcc.EVCC, powerMeter powerMeter.PowerMeter, maxPower float64, maxChargingTime *int) (*Impl, error) {
	log.WithFields(log.Fields{
		"evseId":          evseId,
		"maxChargingTime": maxChargingTime,
	}).Info("Creating a new evse")

	if evseId <= 0 {
		return nil, ErrInvalidEvseId
	}

	if util.IsNilInterfaceOrPointer(evcc) {
		return nil, ErrInvalidEVCC
	}

	return &Impl{
		mu:              sync.Mutex{},
		evseId:          evseId,
		evcc:            evcc,
		powerMeter:      powerMeter,
		maxChargingTime: maxChargingTime,
		maxPower:        maxPower,
		status:          core.ChargePointStatusAvailable,
		scheduler:       scheduler.NewScheduler(),
		logger:          log.StandardLogger().WithField("component", "evse").WithField("evseId", evseId),
	}, nil
}

func (evse *Impl) Init(ctx context.Context) error {
	evse.logger.Info("Initializing evse")
	// Init EVCC
	err := evse.evcc.Init(ctx)
	if err != nil {
		return err
	}

	// Disable charging by default
	evse.evcc.DisableCharging()

	// Set max charging current
	err = evse.evcc.SetMaxChargingCurrent(evse.maxPower)
	if err != nil {
		return err
	}

	// Listen for EVCC status updates in another thread
	go evse.listenForStatusUpdates(ctx)
	return nil
}

func (evse *Impl) listenForStatusUpdates(ctx context.Context) {
	evse.logger.Debug("Listening for evcc status updates")

	statusChan := evse.evcc.GetStatusChangeChannel()
	if statusChan == nil {
		log.Panic("Cannot listen for evcc status updates")
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

			// Compare to current status
			switch evse.status {
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

			evse.SetStatus(state, cpErr)
		case <-ctx.Done():
			break Loop
		}
	}
}

// StartCharging Start charging an evse if evse is available and session could be started.
func (evse *Impl) StartCharging(connectorId *int, measurands []types.Measurand, sampleInterval string) error {
	logInfo := evse.logger.WithField("connectorId", connectorId)
	logInfo.Debugf("Trying to start charging on evse")

	// Check if evse is available
	if !(evse.IsAvailable() || evse.IsPreparing()) {
		return ErrInvalidStatus
	}

	// Enable charging
	err := evse.evcc.EnableCharging()
	if err != nil {
		return err
	}

	evse.evcc.Lock()

	// Prepare power meter and schedule sampling
	sampleError := evse.scheduleMeterValueUpdates(measurands, sampleInterval)
	if sampleError != nil {
		logInfo.WithError(sampleError).Error("Cannot sample evse")
	}

	// Schedule a stop charging after the maxChargingTime, if provided
	if evse.maxChargingTime != nil {
		// Schedule a stop charging after the maxChargingTime
		_, err := evse.scheduler.Every(*evse.maxChargingTime).Minutes().Tag("evse", fmt.Sprintf("%d", evse.GetEvseId()), "chargingTimer").Do(evse.StopCharging, core.ReasonLocal)
		if err != nil {
			logInfo.WithError(err).Error("Cannot schedule stop charging")
		}
	}

	return nil
}

// StopCharging Stops charging an evse if evse is charging
func (evse *Impl) StopCharging(reason core.Reason) error {
	logInfo := evse.logger.WithField("reason", reason)

	if evse.IsCharging() || evse.IsPreparing() {
		logInfo.Debugf("Stopping charging")

		evse.evcc.DisableCharging()
		evse.evcc.Unlock()

		// Remove any jobs scheduled for this evse
		schedulerErr := evse.scheduler.RemoveByTag(fmt.Sprintf("%d", evse.GetEvseId()))
		if schedulerErr != nil {
			logInfo.WithError(schedulerErr).Errorf("Cannot remove sampling schedule")
		}

		return nil
	}

	return ErrNotCharging
}

func (evse *Impl) GetEvseId() int {
	evse.logger.Debugf("Getting evse id")
	return evse.evseId
}

func (evse *Impl) GetMaxChargingTime() *int {
	evse.logger.Debugf("Getting max charging time")
	return evse.maxChargingTime
}

func (evse *Impl) SetMaxChargingTime(time *int) {
	evse.logger.WithField("time", time).Debugf("Setting max charging time")
	evse.maxChargingTime = time
}

func (evse *Impl) GetMaxChargingPower() float64 {
	evse.logger.Debugf("Getting max charging power")
	return evse.maxPower
}

func (evse *Impl) SetNotificationChannel(notificationChannel chan<- notifications.StatusNotification) {
	evse.notificationChannel = notificationChannel
}

func (evse *Impl) SetMeterValuesChannel(notificationChannel chan<- notifications.MeterValueNotification) {
	evse.meterValuesChannel = notificationChannel
}

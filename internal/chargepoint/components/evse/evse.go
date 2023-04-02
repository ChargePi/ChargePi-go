package evse

import (
	"errors"
	"fmt"
	"sync"

	"github.com/go-co-op/gocron"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/evcc"
	"github.com/xBlaz3kx/ChargePi-go/pkg/power-meter"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"golang.org/x/net/context"
)

var (
	ErrInvalidEvseId            = errors.New("invalid evse id")
	ErrInvalidReservationId     = errors.New("invalid reservation id")
	ErrInvalidStatus            = errors.New("invalid evse status")
	ErrInvalidEVCC              = errors.New("evcc cannot be nil")
	ErrSessionTimeLimitExceeded = errors.New("session time limit exceeded")
	ErrNotCharging              = errors.New("evse not charging")
	ErrPowerMeterNotEnabled     = errors.New("power meter not enabled")
	ErrConnectorExists          = errors.New("connector already exists")
)

type (
	EVSE interface {
		Init(ctx context.Context) error

		AddConnector(connector settings.Connector) error
		GetConnectors() []settings.Connector
		GetMaxChargingTime() *int
		SetMaxChargingTime(time *int)
		GetMaxChargingPower() float64

		StartCharging(connectorId *int) error
		StopCharging(reason core.Reason) error
		Lock()
		Unlock()

		GetEvseId() int

		GetPowerMeter() powerMeter.PowerMeter
		SetPowerMeter(powerMeter.PowerMeter) error
		SamplePowerMeter(measurands []types.Measurand)

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
	}, nil
}

func (evse *Impl) Init(ctx context.Context) error {
	// Init EVCC
	err := evse.evcc.Init(ctx)
	if err != nil {
		return err
	}

	// Disable charging by default
	evse.evcc.DisableCharging()

	// Listen for EVCC status updates in another thread
	go evse.listenForStatusUpdates(ctx)
	return nil
}

func (evse *Impl) listenForStatusUpdates(ctx context.Context) {
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

			case core.ChargePointStatusSuspendedEV:
			case core.ChargePointStatusSuspendedEVSE:
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

// StartCharging Start charging a evse if evse is available and session could be started.
// It turns on the relay (even if negative logic applies).
func (evse *Impl) StartCharging(connectorId *int) error {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
	})
	logInfo.Debugf("Trying to start charging on evse")

	if !(evse.IsAvailable() || evse.IsPreparing()) {
		return ErrInvalidStatus
	}

	err := evse.evcc.EnableCharging()
	if err != nil {
		return err
	}

	evse.evcc.Lock()

	sampleError := evse.preparePowerMeter()
	if sampleError != nil {
		logInfo.WithError(sampleError).Error("Cannot sample evse")
	}

	return nil
}

// StopCharging Stops charging the evse by turning the relay off and ending the session.
func (evse *Impl) StopCharging(reason core.Reason) error {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
		"reason": reason,
	})

	if evse.IsCharging() || evse.IsPreparing() {
		logInfo.Debugf("Stopping charging")

		evse.evcc.DisableCharging()
		evse.evcc.Unlock()

		// Remove the sampling of the power meter
		schedulerErr := evse.scheduler.RemoveByTag(fmt.Sprintf("evse%dSampling", evse.GetEvseId()))
		if schedulerErr != nil {
			logInfo.WithError(schedulerErr).Errorf("Cannot remove sampling schedule")
		}

		return nil
	}

	return ErrNotCharging
}

func (evse *Impl) GetEvseId() int {
	return evse.evseId
}

func (evse *Impl) GetMaxChargingTime() *int {
	return evse.maxChargingTime
}

func (evse *Impl) SetMaxChargingTime(time *int) {
	evse.maxChargingTime = time
}

func (evse *Impl) GetMaxChargingPower() float64 {
	return evse.maxPower
}

func (evse *Impl) SetNotificationChannel(notificationChannel chan<- notifications.StatusNotification) {
	evse.notificationChannel = notificationChannel
}

func (evse *Impl) SetMeterValuesChannel(notificationChannel chan<- notifications.MeterValueNotification) {
	evse.meterValuesChannel = notificationChannel
}

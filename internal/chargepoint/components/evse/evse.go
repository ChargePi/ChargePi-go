package evse

import (
	"errors"
	"fmt"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/session"
	"sync"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/evcc"
	powerMeter "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	carState "github.com/xBlaz3kx/ChargePi-go/pkg/models/evcc"
	"golang.org/x/net/context"
)

var (
	ErrInvalidEvseId            = errors.New("invalid evse id")
	ErrInvalidReservationId     = errors.New("invalid reservation id")
	ErrInvalidStatus            = errors.New("invalid evse status")
	ErrInvalidEVCC              = errors.New("evcc cannot be nil")
	ErrSessionTimeLimitExceeded = errors.New("session time limit exceeded")
	ErrNotCharging              = errors.New("evse not charging")
)

type (
	EVSE interface {
		Init(ctx context.Context) error

		AddConnector(connector Connector) error
		GetConnectors() []Connector
		GetMaxChargingTime() *int
		SetMaxChargingTime(time *int)
		GetMaxChargingPower() float64

		StartCharging(transactionId, tagId string, connectorId *int) error
		ResumeCharging(session session.Session) (chargingTimeElapsed *int, err error)
		StopCharging(reason core.Reason) error
		Lock()
		Unlock()

		GetSession() session.Session
		GetTagId() string
		GetTransactionId() string
		GetEvseId() int

		GetPowerMeter() powerMeter.PowerMeter
		SetPowerMeter(powerMeter.PowerMeter) error
		SamplePowerMeter(measurands []types.Measurand)
		CalculateSessionAvgEnergyConsumption() float64

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

		Reserve(reservationId int, tagId string) error
		RemoveReservation() error
		GetReservationId() int
	}

	Impl struct {
		evseId          int
		maxPower        float64
		maxChargingTime *int
		connectors      []Connector
		availability    core.AvailabilityType
		status          core.ChargePointStatus
		errorCode       core.ChargePointErrorCode
		session         *session.Session
		reservationId   *int

		// Notification channels
		meterValuesChannel  chan<- notifications.MeterValueNotification
		notificationChannel chan<- notifications.StatusNotification
		mu                  sync.Mutex

		// Hardware
		powerMeterEnabled bool
		powerMeter        powerMeter.PowerMeter
		evcc              evcc.EVCC
	}
)

// NewEvse Create a new evse object from the provided arguments. evseId, connectorId and maxChargingTime must be greater than zero.
// When created, it makes an empty session, turns off the relay and defaults the status to Available.
func NewEvse(evseId int, evcc evcc.EVCC, powerMeter powerMeter.PowerMeter, powerMeterEnabled bool, maxPower float64, maxChargingTime *int) (*Impl, error) {
	log.WithFields(log.Fields{
		"evseId":          evseId,
		"maxChargingTime": maxChargingTime,
		"hasPowerMeter":   powerMeterEnabled,
	}).Info("Creating a new evse")

	if evseId <= 0 {
		return nil, ErrInvalidEvseId
	}

	if util.IsNilInterfaceOrPointer(evcc) {
		return nil, ErrInvalidEVCC
	}

	return &Impl{
		mu:                sync.Mutex{},
		evseId:            evseId,
		evcc:              evcc,
		powerMeter:        powerMeter,
		powerMeterEnabled: powerMeterEnabled,
		maxChargingTime:   maxChargingTime,
		maxPower:          maxPower,
		status:            core.ChargePointStatusAvailable,
		session:           session.NewEmptySession(),
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
				case carState.StateB1:
					state = core.ChargePointStatusPreparing
				}

			case core.ChargePointStatusPreparing:

				// Determine new state based on the previous state
				switch msg.State {
				case carState.StateA1, carState.StateA2:
					state = core.ChargePointStatusAvailable
				case carState.StateC2, carState.StateD2:
					state = core.ChargePointStatusCharging
				}

			case core.ChargePointStatusCharging:

				switch msg.State {
				case carState.StateC1:
					state = core.ChargePointStatusFinishing
				case carState.StateD1:
					state = core.ChargePointStatusSuspendedEV
				}

			case core.ChargePointStatusSuspendedEV:
			case core.ChargePointStatusSuspendedEVSE:
			case core.ChargePointStatusFaulted:
			}

			switch msg.State {
			case carState.StateE, carState.StateF:
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
func (evse *Impl) StartCharging(transactionId, tagId string, connectorId *int) error {
	logInfo := log.WithFields(log.Fields{
		"evseId":        evse.evseId,
		"transactionId": transactionId,
		"tagId":         tagId,
	})
	logInfo.Debugf("Trying to start charging on evse")

	if !(evse.IsAvailable() || evse.IsPreparing()) {
		return ErrInvalidStatus
	}

	sessionErr := evse.session.StartSession(transactionId, tagId)
	if sessionErr != nil {
		return sessionErr
	}

	err := evse.evcc.EnableCharging()
	if err != nil {
		return err
	}

	evse.evcc.Lock()
	evse.session.UpdateSessionFile(evse.evseId)

	sampleError := evse.preparePowerMeterAtConnector()
	if sampleError != nil {
		logInfo.WithError(sampleError).Error("Cannot sample evse")
	}

	return nil
}

// ResumeCharging Resumes or restores the charging state after boot if a charging session was active.
func (evse *Impl) ResumeCharging(session session.Session) (chargingTimeElapsed *int, err error) {
	logInfo := log.WithFields(log.Fields{
		"evseId":  evse.evseId,
		"session": session,
	})
	logInfo.Debugf("Trying to resume charging on evse")

	// Set the transaction id so evse is able to stop the transaction if charging fails
	evse.session.TransactionId = session.TransactionId

	sessionErr := evse.session.StartSession(session.TransactionId, session.TagId)
	if sessionErr != nil {
		return evse.maxChargingTime, fmt.Errorf("cannot resume session: %v", sessionErr)
	}

	if evse.IsPreparing() || evse.IsCharging() {
		err = evse.evcc.EnableCharging()
		if err != nil {
			return
		}

		evse.session.Started = session.Started
		evse.session.Consumption = session.Consumption
	}

	startedChargingTime, err := time.Parse(time.RFC3339, session.Started)
	if err != nil {
		return
	}

	if evse.maxChargingTime != nil {
		timeElapsed := int(time.Now().Sub(startedChargingTime).Minutes())
		if *evse.maxChargingTime <= timeElapsed {
			chargingTimeElapsed = &timeElapsed
		}

		err = ErrSessionTimeLimitExceeded
		return
	}

	return nil, nil
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
		evse.session.EndSession()
		evse.session.UpdateSessionFile(evse.evseId)

		// Remove the sampling of the power meter
		sched := scheduler.GetScheduler()
		schedulerErr := sched.RemoveByTag(fmt.Sprintf("evse%dSampling", evse.GetEvseId()))
		if schedulerErr != nil {
			logInfo.WithError(schedulerErr).Errorf("Cannot remove sampling schedule")
		}

		return nil
	}

	return ErrNotCharging
}

func (evse *Impl) GetSession() session.Session {
	return *evse.session
}

func (evse *Impl) GetTransactionId() string {
	return evse.session.TransactionId
}

func (evse *Impl) GetTagId() string {
	return evse.session.TagId
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

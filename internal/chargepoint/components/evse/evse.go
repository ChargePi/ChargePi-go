package evse

import (
	"errors"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/evcc"
	powerMeter "github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/internal/models"
	carState "github.com/xBlaz3kx/ChargePi-go/internal/models/evcc"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/session"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"golang.org/x/net/context"
	"sync"
	"time"
)

var (
	ErrInvalidEvseId            = errors.New("invalid evse id")
	ErrInvalidReservationId     = errors.New("invalid reservation id")
	ErrInvalidStatus            = errors.New("invalid evse status")
	ErrRelayPointerNil          = errors.New("relay pointer cannot be nil")
	ErrSessionTimeLimitExceeded = errors.New("session time limit exceeded")
	ErrNotCharging              = errors.New("evse not charging")
)

type (
	EVSE interface {
		Init(ctx context.Context) error
		StartCharging(transactionId, tagId string, connectorId *int) error
		ResumeCharging(session session.Session) (int, error)
		StopCharging(reason core.Reason) error
		GetTagId() string
		GetTransactionId() string
		GetEvseId() int
		GetConnectors() []Connector
		SetAvailability(isAvailable bool)

		SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode)
		GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode)
		IsAvailable() bool
		IsPreparing() bool
		IsCharging() bool
		IsReserved() bool
		IsUnavailable() bool
		GetMaxChargingTime() int

		SetNotificationChannel(notificationChannel chan<- models.StatusNotification)
		SetMeterValuesChannel(notificationChannel chan<- models.MeterValueNotification)

		ReserveEvse(reservationId int, tagId string) error
		RemoveReservation() error
		GetReservationId() int

		CalculateSessionAvgEnergyConsumption() float64
		SamplePowerMeter(measurands []types.Measurand)
	}

	evseImpl struct {
		evseId          int
		maxChargingTime int
		availability    core.AvailabilityType
		status          core.ChargePointStatus
		errorCode       core.ChargePointErrorCode
		session         *session.Session
		reservationId   *int

		// Notification channels
		meterValuesChannel  chan<- models.MeterValueNotification
		notificationChannel chan<- models.StatusNotification
		mu                  sync.Mutex

		// Hardware
		powerMeter        powerMeter.PowerMeter
		powerMeterEnabled bool
		evcc              evcc.EVCC
	}
)

// NewEvse Create a new evse object from the provided arguments. evseId, connectorId and maxChargingTime must be greater than zero.
// When created, it makes an empty session, turns off the relay and defaults the status to Available.
func NewEvse(evseId int, evcc evcc.EVCC, powerMeter powerMeter.PowerMeter, powerMeterEnabled bool, maxChargingTime int) (*evseImpl, error) {
	log.WithFields(log.Fields{
		"evseId":          evseId,
		"maxChargingTime": maxChargingTime,
		"hasPowerMeter":   powerMeterEnabled,
	}).Info("Creating a new evse")

	if maxChargingTime <= 0 {
		maxChargingTime = 180
	}

	if evseId <= 0 {
		return nil, ErrInvalidEvseId
	}

	if util.IsNilInterfaceOrPointer(evcc) {
		return nil, ErrRelayPointerNil
	}

	return &evseImpl{
		mu:                sync.Mutex{},
		evseId:            evseId,
		evcc:              evcc,
		powerMeter:        powerMeter,
		powerMeterEnabled: powerMeterEnabled,
		maxChargingTime:   maxChargingTime,
		status:            core.ChargePointStatusAvailable,
		session:           session.NewEmptySession(),
	}, nil
}

func (evse *evseImpl) Init(ctx context.Context) error {
	// Init EVCC
	err := evse.evcc.Init(ctx)
	if err != nil {
		return err
	}

	// Disable charging
	evse.evcc.DisableCharging()

	statusChan := evse.evcc.GetStatusChangeChannel()
	if statusChan == nil {
		return nil
	}

	// Listen for EVCC status updates in another thread
	go func() {

	Loop:
		for {
			select {
			case msg := <-statusChan:
				// Determine OCPP status based on CarState and Error

				var (
					state core.ChargePointStatus
					cpErr core.ChargePointErrorCode
				)

				switch msg.State {
				case carState.StateA1:
					state = core.ChargePointStatusAvailable

				default:
					state = core.ChargePointStatusFaulted
				}

				evse.SetStatus(state, cpErr)
			case <-ctx.Done():
				break Loop
			}
		}
	}()

	return nil
}

// StartCharging Start charging a evse if evse is available and session could be started.
// It turns on the relay (even if negative logic applies).
func (evse *evseImpl) StartCharging(transactionId, tagId string, connectorId *int) error {
	logInfo := log.WithFields(log.Fields{
		"evseId":        evse.evseId,
		"transactionId": transactionId,
		"tagId":         tagId,
	})
	logInfo.Debugf("Trying to start charging on evse")

	if !(evse.IsAvailable() || evse.IsPreparing()) {
		return ErrInvalidStatus
	}

	//evse.SetStatus(core.ChargePointStatusPreparing, core.NoError)
	sessionErr := evse.session.StartSession(transactionId, tagId)
	if sessionErr != nil {
		return sessionErr
	}

	err := evse.evcc.EnableCharging()
	if err != nil {
		return err
	}

	evse.session.UpdateSessionFile(evse.evseId)

	if evse.powerMeterEnabled && !util.IsNilInterfaceOrPointer(evse.powerMeter) {
		sampleError := evse.preparePowerMeterAtConnector()
		if sampleError != nil {
			logInfo.Errorf("Cannot sample evse: %v", sampleError)
		}
	}

	return nil
}

// ResumeCharging Resumes or restores the charging state after boot if a charging session was active.
func (evse *evseImpl) ResumeCharging(session session.Session) (chargingTimeElapsed int, err error) {
	// Set the transaction id so evse is able to stop the transaction if charging fails
	logInfo := log.WithFields(log.Fields{
		"evseId":  evse.evseId,
		"session": session,
	})
	logInfo.Debugf("Trying to resume charging on evse")

	chargingTimeElapsed = evse.maxChargingTime
	evse.session.TransactionId = session.TransactionId

	startedChargingTime, err := time.Parse(time.RFC3339, session.Started)
	if err != nil {
		return
	}

	chargingTimeElapsed = int(time.Now().Sub(startedChargingTime).Minutes())
	if evse.maxChargingTime <= chargingTimeElapsed {
		chargingTimeElapsed = evse.maxChargingTime
		err = ErrSessionTimeLimitExceeded
		return
	}

	if evse.IsCharging() || evse.IsPreparing() {
		sessionErr := evse.session.StartSession(session.TransactionId, session.TagId)
		if sessionErr != nil {
			return evse.maxChargingTime, fmt.Errorf("cannot resume session: %v", sessionErr)
		}

		err = evse.evcc.EnableCharging()
		if err != nil {
			return
		}

		evse.session.Started = session.Started
		evse.session.Consumption = append(evse.session.Consumption, session.Consumption...)
		return chargingTimeElapsed, nil
	}

	return evse.maxChargingTime, ErrInvalidStatus
}

// StopCharging Stops charging the evse by turning the relay off and ending the session.
func (evse *evseImpl) StopCharging(reason core.Reason) error {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
		"reason": reason,
	})

	if evse.IsCharging() || evse.IsPreparing() {
		logInfo.Debugf("Stopping charging")

		evse.evcc.DisableCharging()
		evse.session.EndSession()
		evse.session.UpdateSessionFile(evse.evseId)
		return nil
	}

	return ErrNotCharging
}

// SamplePowerMeter Get a sample from the power meter. The measurands argument takes the list of all the types of the measurands to sample.
// It will add all the samples to the evse's Session if it is active.
func (evse *evseImpl) SamplePowerMeter(measurands []types.Measurand) {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
	})

	if !evse.powerMeterEnabled || util.IsNilInterfaceOrPointer(evse.powerMeter) {
		return
	}

	logInfo.Debugf("Sampling evse %v", measurands)
	var (
		meterValues []types.MeterValue
		samples     []types.SampledValue
		value       = 0.0
	)

	for _, measurand := range measurands {

		switch measurand {
		case types.MeasurandEnergyActiveImportInterval, types.MeasurandEnergyActiveImportRegister,
			types.MeasurandEnergyActiveExportInterval, types.MeasurandEnergyActiveExportRegister:
			value = evse.powerMeter.GetEnergy()
		case types.MeasurandCurrentImport, types.MeasurandCurrentExport:
			value = evse.powerMeter.GetCurrent()
		case types.MeasurandPowerActiveImport, types.MeasurandPowerActiveExport:
			value = evse.powerMeter.GetPower()
		case types.MeasurandVoltage:
			value = evse.powerMeter.GetVoltage()
		}

		if value != 0.0 {
			sample := types.SampledValue{
				Value:     fmt.Sprintf("%.3f", value),
				Measurand: measurand,
			}

			meterValues = append(meterValues, types.MeterValue{SampledValue: []types.SampledValue{sample}, Timestamp: types.NewDateTime(time.Now())})
		}
	}

	if evse.meterValuesChannel != nil {
		evse.meterValuesChannel <- models.NewMeterValueNotification(evse.evseId, nil, nil, meterValues...)
	}

	evse.session.AddSampledValue(samples)
}

// preparePowerMeterAtConnector
func (evse *evseImpl) preparePowerMeterAtConnector() error {
	var (
		measurands          = util.GetTypesToSample()
		sampleTime          = "10s"
		sampleInterval, err = ocppConfigManager.GetConfigurationValue(v16.MeterValueSampleInterval.String())
		jobTag              = fmt.Sprintf("Evse%dSampling", evse.evseId)
	)
	if err != nil {
		sampleInterval = "10"
	}

	sampleTime = fmt.Sprintf("%ss", sampleInterval)
	// Schedule the sampling
	_, err = scheduler.GetScheduler().Every(sampleTime).
		Tag(jobTag).
		Do(evse.SamplePowerMeter, measurands)

	return err
}

func (evse *evseImpl) IsAvailable() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusAvailable && evse.availability == core.AvailabilityTypeOperative
}

func (evse *evseImpl) IsCharging() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusCharging
}

func (evse *evseImpl) IsPreparing() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusPreparing
}

func (evse *evseImpl) IsReserved() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusReserved
}

func (evse *evseImpl) IsUnavailable() bool {
	evse.mu.Lock()
	defer evse.mu.Unlock()
	return evse.status == core.ChargePointStatusUnavailable
}

func (evse *evseImpl) SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
	})
	logInfo.Debugf("Setting evse status %s with err %s", status, errCode)

	evse.mu.Lock()
	defer evse.mu.Unlock()

	evse.status = status
	evse.errorCode = errCode
	settings.UpdateEVSEStatus(evse.evseId, status)

	if evse.notificationChannel != nil {
		evse.notificationChannel <- models.NewStatusNotification(evse.evseId, string(status), string(errCode))
	}
}

func (evse *evseImpl) GetTransactionId() string {
	return evse.session.TransactionId
}

func (evse *evseImpl) GetTagId() string {
	return evse.session.TagId
}

func (evse *evseImpl) ReserveEvse(reservationId int, tagId string) error {
	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
		"tagId":  tagId,
	})
	logInfo.Debugf("Reserving evse for id %d", reservationId)

	if reservationId <= 0 {
		return ErrInvalidReservationId
	}

	if !evse.IsAvailable() {
		return ErrInvalidStatus
	}

	evse.reservationId = &reservationId
	evse.SetStatus(core.ChargePointStatusReserved, core.NoError)
	return nil
}

func (evse *evseImpl) RemoveReservation() error {
	if !evse.IsReserved() {
		return ErrInvalidStatus
	}

	logInfo := log.WithFields(log.Fields{
		"evseId": evse.evseId,
	})
	logInfo.Debugf("Removing reservation")

	evse.reservationId = nil
	evse.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	return nil
}

func (evse *evseImpl) GetReservationId() int {
	if util.IsNilInterfaceOrPointer(evse.reservationId) {
		return -1
	}

	return *evse.reservationId
}

func (evse *evseImpl) GetEvseId() int {
	return evse.evseId
}

func (evse *evseImpl) CalculateSessionAvgEnergyConsumption() float64 {
	return evse.session.CalculateEnergyConsumptionWithAvgPower()
}

func (evse *evseImpl) GetMaxChargingTime() int {
	return evse.maxChargingTime
}

func (evse *evseImpl) GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode) {
	return evse.status, evse.errorCode
}

func (evse *evseImpl) SetNotificationChannel(notificationChannel chan<- models.StatusNotification) {
	evse.notificationChannel = notificationChannel
}

func (evse *evseImpl) SetMeterValuesChannel(notificationChannel chan<- models.MeterValueNotification) {
	evse.meterValuesChannel = notificationChannel
}

func (evse *evseImpl) GetConnectors() []Connector {
	return nil
}

func (evse *evseImpl) SetAvailability(isAvailable bool) {
	if isAvailable {
		evse.availability = core.AvailabilityTypeOperative
		return
	}

	evse.availability = core.AvailabilityTypeInoperative
}

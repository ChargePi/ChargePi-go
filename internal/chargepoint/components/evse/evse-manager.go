package evse

import (
	"context"
	"errors"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/evcc"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/power-meter"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/session"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"sync"
)

var (
	ErrConnectorNotFound      = errors.New("connector not found")
	ErrConnectorAlreadyExists = errors.New("connector already exists")
	ErrConnectorStatusInvalid = errors.New("connector status invalid")
	ErrConnectorNil           = errors.New("connector is nil")

	manager Manager
	once    = sync.Once{}
)

type (
	Manager interface {
		GetEVSEs() []EVSE
		FindEVSE(evseId int) (EVSE, error)
		FindAvailableEVSE() (EVSE, error)
		FindEVSEWithTagId(tagId string) (EVSE, error)
		FindEVSEWithTransactionId(transactionId string) (EVSE, error)
		FindEVSEWithReservationId(reservationId int) (EVSE, error)
		StartCharging(evseId int, tagId, transactionId string) error
		StopCharging(tagId, transactionId string, reason core.Reason) error
		StopAllEVSEs(reason core.Reason) error
		AddEVSE(c EVSE) error
		AddEVSEFromSettings(maxChargingTime *int, c *settings.EVSE) error
		AddEVSEsFromSettings(maxChargingTime *int, c []*settings.EVSE) error
		RestoreEVSEStatus(*settings.EVSE) error
		SetNotificationChannel(notificationChannel chan chargePoint.StatusNotification)
		SetMeterValuesChannel(notificationChannel chan chargePoint.MeterValueNotification)
	}

	managerImpl struct {
		connectors          sync.Map
		notificationChannel chan chargePoint.StatusNotification
		meterValuesChannel  chan chargePoint.MeterValueNotification
	}
)

func init() {
	once.Do(func() {
		GetManager()
	})
}

func GetManager() Manager {
	if manager == nil {
		log.Debug("Creating connector manager")
		manager = NewManager(nil)
	}

	return manager
}

func NewManager(notificationChannel chan chargePoint.StatusNotification) Manager {
	return &managerImpl{
		connectors:          sync.Map{},
		notificationChannel: notificationChannel,
	}
}

func (m *managerImpl) GetEVSEs() []EVSE {
	var connectors []EVSE

	m.connectors.Range(func(key, value interface{}) bool {
		c, canCast := value.(EVSE)
		if canCast {
			connectors = append(connectors, c)
		}
		return true
	})

	return connectors
}

func (m *managerImpl) SetNotificationChannel(notificationChannel chan chargePoint.StatusNotification) {
	if notificationChannel != nil {
		m.notificationChannel = notificationChannel
	}
}

func (m *managerImpl) SetMeterValuesChannel(notificationChannel chan chargePoint.MeterValueNotification) {
	if notificationChannel != nil {
		m.meterValuesChannel = notificationChannel
	}
}

func (m *managerImpl) FindEVSE(evseId int) (EVSE, error) {
	var (
		key        = fmt.Sprintf("Evse%d", evseId)
		c, isFound = m.connectors.Load(key)
	)

	if isFound {
		return c.(EVSE), nil
	}

	return nil, ErrConnectorNotFound
}

func (m *managerImpl) FindAvailableEVSE() (EVSE, error) {
	var availableConnector EVSE

	m.connectors.Range(func(key, value interface{}) bool {
		c, canCast := value.(EVSE)
		if canCast && c.IsAvailable() {
			availableConnector = c
			return false
		}

		return true
	})

	if util.IsNilInterfaceOrPointer(availableConnector) {
		return nil, ErrConnectorNotFound
	}

	return availableConnector, nil
}

func (m *managerImpl) FindEVSEWithTagId(tagId string) (EVSE, error) {
	var connectorWithTag EVSE

	m.connectors.Range(func(key, value interface{}) bool {
		c, _ := value.(EVSE)
		if c.GetTagId() == tagId {
			connectorWithTag = c
			return false
		}

		return true
	})

	if util.IsNilInterfaceOrPointer(connectorWithTag) {
		return nil, ErrConnectorNotFound
	}

	return connectorWithTag, nil
}

func (m *managerImpl) FindEVSEWithTransactionId(transactionId string) (EVSE, error) {
	var connectorWithTransaction EVSE

	m.connectors.Range(func(key, value interface{}) bool {
		c, _ := value.(EVSE)
		if c.GetTransactionId() == transactionId {
			connectorWithTransaction = c
			return false
		}

		return true
	})

	if util.IsNilInterfaceOrPointer(connectorWithTransaction) {
		return nil, ErrConnectorNotFound
	}

	return connectorWithTransaction, nil
}

func (m *managerImpl) FindEVSEWithReservationId(reservationId int) (EVSE, error) {
	var connectorWithReservation EVSE

	m.connectors.Range(func(key, value interface{}) bool {
		c, _ := value.(EVSE)
		if c.GetReservationId() == reservationId {
			connectorWithReservation = c
			return false
		}

		return true
	})

	if util.IsNilInterfaceOrPointer(connectorWithReservation) {
		return nil, ErrConnectorNotFound
	}

	return connectorWithReservation, nil
}

func (m *managerImpl) StartCharging(evseId int, tagId, transactionId string) error {
	c, _ := m.FindEVSE(evseId)

	if c != nil {
		return c.StartCharging(transactionId, tagId, nil)
	}

	return ErrConnectorNotFound
}

func (m *managerImpl) StopCharging(tagId, transactionId string, reason core.Reason) error {
	c, _ := m.FindEVSEWithTransactionId(transactionId)

	if c != nil {
		return c.StopCharging(reason)
	}

	return ErrConnectorNotFound
}

func (m *managerImpl) StopAllEVSEs(reason core.Reason) error {
	log.Debugf("Stopping all evses: %s", reason)

	var err error

	for _, c := range m.GetEVSEs() {
		stopErr := c.StopCharging(reason)
		if stopErr != nil {
			err = stopErr
		}
	}

	return err
}

func (m *managerImpl) AddEVSE(c EVSE) error {
	if util.IsNilInterfaceOrPointer(c) {
		return ErrConnectorNil
	}

	var (
		logInfo = log.WithFields(log.Fields{
			"evseId": c.GetEvseId(),
		})
		key = fmt.Sprintf("Evse%d", c.GetEvseId())
		ctx = context.Background()
	)

	err := c.Init(ctx)
	if err != nil {
		return err
	}

	logInfo.Debugf("Adding an EVSE to manager")
	c.SetNotificationChannel(m.notificationChannel)
	c.SetMeterValuesChannel(m.meterValuesChannel)

	// Add the connector
	m.connectors.Store(key, c)

	return nil
}

func (m *managerImpl) AddEVSEFromSettings(maxChargingTime *int, c *settings.EVSE) error {
	if util.IsNilInterfaceOrPointer(c) {
		return ErrConnectorNil
	}

	var (
		evccFromType, err = evcc.NewEVCCFromType(c.EVCC)
		// Create a PowerMeter from connector settings
		meter, powerMeterErr = powerMeter.NewPowerMeter(c.PowerMeter)
	)

	switch powerMeterErr {
	case powerMeter.ErrPowerMeterDisabled:
		log.WithError(powerMeterErr).Warnf("Power meter disabled for evse %d", c.EvseId)
	case powerMeter.ErrPowerMeterUnsupported, powerMeter.ErrInvalidConnectionSettings:
		fallthrough
	default:
		log.WithError(powerMeterErr).Fatalf("Cannot instantiate power meter for evse %d", c.EvseId)
	}

	evse, err := NewEvse(c.EvseId, evccFromType, meter, c.PowerMeter.Enabled, maxChargingTime)
	if err != nil {
		return err
	}

	return m.AddEVSE(evse)
}

func (m *managerImpl) AddEVSEsFromSettings(maxChargingTime *int, connectors []*settings.EVSE) error {
	for _, c := range connectors {
		addErr := m.AddEVSEFromSettings(maxChargingTime, c)
		if addErr != nil {
			return addErr
		}
	}

	return nil
}

func (m *managerImpl) RestoreEVSEStatus(c *settings.EVSE) error {
	if util.IsNilInterfaceOrPointer(c) {
		return ErrConnectorNotFound
	}

	logInfo := log.WithFields(log.Fields{
		"evseId":  c.EvseId,
		"session": c.Session,
	})
	logInfo.Debugf("Attempting to restore connector status")

	// Current status
	evse, err := m.FindEVSE(c.EvseId)
	if err != nil {
		return err
	}

	// Determine what to do based on the previous status
	switch core.ChargePointStatus(c.Status) {
	case core.ChargePointStatusAvailable:
		return nil
	case core.ChargePointStatusPreparing:
		return evse.StartCharging(c.Session.TransactionId, c.Session.TagId, nil)
	case core.ChargePointStatusCharging:
		timeLeft, err := evse.ResumeCharging(session.Session(c.Session))
		if err != nil {
			logInfo.Errorf("Resume charging failed, reason: %v", err)

			// Attempt to stop charging
			err = evse.StopCharging(core.ReasonDeAuthorized)
			if err != nil {
				logInfo.Errorf("Err stopping the charging: %v", err)
				evse.SetStatus(core.ChargePointStatusFaulted, core.InternalError)
			}
		}

		// Stop charging
		_, schedulerErr := scheduler.GetScheduler().Every(*timeLeft).Minutes().At(1).Do(evse.StopCharging(core.ReasonLocal))
		if schedulerErr != nil {
			return schedulerErr
		}

		logInfo.Debugf("Successfully resumed charging")
		return nil
	case core.ChargePointStatusReserved:
		//todo
		return nil
	case core.ChargePointStatusFaulted,
		core.ChargePointStatusFinishing,
		core.ChargePointStatusSuspendedEV,
		core.ChargePointStatusSuspendedEVSE:
		return nil
	default:
		return ErrConnectorStatusInvalid
	}
}

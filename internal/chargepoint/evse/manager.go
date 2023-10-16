package evse

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/evcc"
	"github.com/xBlaz3kx/ChargePi-go/pkg/power-meter"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

var (
	ErrConnectorNotFound      = errors.New("connector not found")
	ErrReservationNotFound    = errors.New("reservation not found")
	ErrConnectorStatusInvalid = errors.New("connector status invalid")
	ErrConnectorNil           = errors.New("connector is nil")

	manager Manager
	once    = sync.Once{}
)

type (
	Manager interface {
		InitAll(ctx context.Context) error
		AddEVSE(ctx context.Context, c EVSE) error
		UpdateEVSE(ctx context.Context, c EVSE) error
		RemoveEVSE(evseId int) error
		GetEVSEs() []EVSE
		GetEVSE(evseId int) (EVSE, error)
		GetAvailableEVSE() (EVSE, error)
		GetEVSEWithReservationId(reservationId int) (EVSE, error)

		StartCharging(evseId int, connectorId *int) error
		StopCharging(evseId int, connectorId *int, reason core.Reason) error
		StopAllEVSEs(reason core.Reason) error
		RestoreEVSEs() error
		// GetCurrentConsumption() (*types.MeterValue,error)
		// UnlockConnector(evseId, connectorId int) error
		// RestoreEVSE(evseId int,) error

		Reserve(evseId int, connectorId *int, reservationId int, tagId string) error
		RemoveReservation(reservationId int) error

		SetNotificationChannel(notificationChannel chan notifications.StatusNotification)
		GetNotificationChannel() chan notifications.StatusNotification
		SetMeterValuesChannel(notificationChannel chan notifications.MeterValueNotification)
	}

	managerImpl struct {
		// Used to store settings of EVSEs
		db *badger.DB
		// Used for running EVSE instances
		connectors          sync.Map
		reservations        map[int]*int
		notificationChannel chan notifications.StatusNotification
		meterValuesChannel  chan notifications.MeterValueNotification
	}
)

func init() {
	once.Do(func() {
		GetManager()
	})
}

func GetManager() Manager {
	if manager == nil {
		log.Debug("Creating EVSE manager")
		manager = NewManager(make(chan notifications.StatusNotification, 20))
	}

	return manager
}

func NewManager(notificationChannel chan notifications.StatusNotification) Manager {
	return &managerImpl{
		db:                  database.Get(),
		notificationChannel: notificationChannel,
	}
}

func getKey(evseId int) string {
	return fmt.Sprintf("evse-%d", evseId)
}

func (m *managerImpl) InitAll(ctx context.Context) error {
	log.Info("Initializing EVSEs")

	// Create EVSEs from settings stored in the database.
	for _, c := range database.GetEvseSettings(m.db) {
		addErr := m.addEVSEFromSettings(ctx, c)
		if addErr != nil {
			return addErr
		}
	}

	return nil
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

func (m *managerImpl) SetNotificationChannel(notificationChannel chan notifications.StatusNotification) {
	if notificationChannel != nil {
		m.notificationChannel = notificationChannel
	}
}

func (m *managerImpl) GetNotificationChannel() chan notifications.StatusNotification {
	return m.notificationChannel
}

func (m *managerImpl) SetMeterValuesChannel(notificationChannel chan notifications.MeterValueNotification) {
	if notificationChannel != nil {
		m.meterValuesChannel = notificationChannel
	}
}

func (m *managerImpl) GetEVSE(evseId int) (EVSE, error) {
	c, isFound := m.connectors.Load(getKey(evseId))
	if isFound {
		return c.(EVSE), nil
	}

	return nil, ErrConnectorNotFound
}

func (m *managerImpl) GetAvailableEVSE() (EVSE, error) {
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

func (m *managerImpl) StartCharging(evseId int, connectorId *int) error {
	c, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	return c.StartCharging(connectorId)
}

func (m *managerImpl) StopCharging(evseId int, connectorId *int, reason core.Reason) error {
	c, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	return c.StopCharging(reason)
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

func (m *managerImpl) AddEVSE(ctx context.Context, c EVSE) error {
	if util.IsNilInterfaceOrPointer(c) {
		return ErrConnectorNil
	}

	err := c.Init(ctx)
	if err != nil {
		return err
	}

	logInfo := log.WithField("evseId", c.GetEvseId())
	logInfo.Debugf("Adding an EVSE to manager")

	c.SetNotificationChannel(m.notificationChannel)
	c.SetMeterValuesChannel(m.meterValuesChannel)

	// Add the connector
	m.connectors.Store(getKey(c.GetEvseId()), c)
	return nil
}

func (m *managerImpl) addEVSEFromSettings(ctx context.Context, c settings.EVSE) error {
	logInfo := log.WithField("evseId", c.EvseId)
	logInfo.Debugf("Creating evcc")

	// Create EVSE from settings
	evccFromType, err := evcc.NewEVCCFromType(c.EVCC)
	switch err {
	case nil:
		logInfo.WithField("type", c.EVCC.Type).Debugf("EVCC created")
	default:
		return err
	}

	// Create a PowerMeter from settings
	logInfo.Debugf("Creating power meter")
	meter, powerMeterErr := powerMeter.NewPowerMeter(c.PowerMeter)
	switch powerMeterErr {
	case nil:
	case powerMeter.ErrPowerMeterDisabled:
		logInfo.WithError(powerMeterErr).Warn("Power meter disabled")
	case powerMeter.ErrPowerMeterUnsupported, powerMeter.ErrInvalidConnectionSettings:
		fallthrough
	default:
		logInfo.WithError(powerMeterErr).Error("Cannot instantiate power meter for evse")
		return err
	}

	// Create EVSE from EVCC and Power Meter
	logInfo.Debugf("Creating EVSE")
	evse, err := NewEvse(c.EvseId, evccFromType, meter, float64(c.MaxPower), nil)
	if err != nil {
		return err
	}

	if m.notificationChannel != nil {
		evse.SetNotificationChannel(m.notificationChannel)
	}

	return m.AddEVSE(ctx, evse)
}

func (m *managerImpl) UpdateEVSE(ctx context.Context, c EVSE) error {
	// todo implement me
	return nil
}

func (m *managerImpl) RemoveEVSE(evseId int) error {
	m.connectors.Delete(getKey(evseId))
	return nil
}

func (m *managerImpl) RestoreEVSEs() error {
	log.Debugf("Attempting to restore EVSEs")

	for _, s := range database.GetEvseSettings(m.db) {
		err := m.restoreEVSEStatus(s)
		if err != nil {
			log.WithError(err).WithField("id", s.EvseId).Error("Error restoring an EVSE")
			continue
		}
	}

	return nil
}

func (m *managerImpl) restoreEVSEStatus(c settings.EVSE) error {
	logInfo := log.WithFields(log.Fields{
		"evseId": c.EvseId,
	})
	logInfo.Debugf("Attempting to restore connector status")

	// Find the EVSE
	evse, err := m.GetEVSE(c.EvseId)
	if err != nil {
		return err
	}

	// Get the current status
	status, _ := evse.GetStatus()

	// Determine what to do based on the previous status
	switch status {
	case core.ChargePointStatusAvailable:
		return nil
	case core.ChargePointStatusPreparing:
		// return evse.StartCharging(c.Session.TransactionId, c.Session.TagId, nil)
		return nil
	case core.ChargePointStatusCharging:

		// Stop charging
		_, schedulerErr := scheduler.NewScheduler().Every(1).Minutes().At(1).Do(evse.StopCharging(core.ReasonLocal))
		if schedulerErr != nil {
			return schedulerErr
		}

		logInfo.Debugf("Successfully resumed charging")
		return nil
	case core.ChargePointStatusReserved:
		// todo
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

package manager

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/pkg/database"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/notifications"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/ChargePi/ChargePi-go/internal/pkg/scheduler"
	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
	"github.com/ChargePi/ChargePi-go/pkg/evcc"
	powerMeter "github.com/ChargePi/ChargePi-go/pkg/power-meter"
	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"go.uber.org/zap"
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
		AddEVSE(ctx context.Context, c evse.EVSE) error
		UpdateEVSE(ctx context.Context, c evse.EVSE) error
		RemoveEVSE(evseId int) error
		GetEVSEs() []evse.EVSE
		GetEVSE(evseId int) (evse.EVSE, error)
		GetAvailableEVSE() (evse.EVSE, error)
		GetEVSEWithReservationId(reservationId int) (evse.EVSE, error)

		StartCharging(evseId int, connectorId *int, measurands []types.Measurand, sampleInterval string) error
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
		logger              *zap.Logger
	}
)

func init() {
	once.Do(func() {
		GetManager()
	})
}

func GetManager() Manager {
	if manager == nil {
		zap.L().Debug("Creating EVSE manager")
		manager = NewManager(make(chan notifications.StatusNotification, 20))
	}

	return manager
}

func NewManager(notificationChannel chan notifications.StatusNotification) Manager {
	return &managerImpl{
		db:                  database.Get(),
		notificationChannel: notificationChannel,
		logger:              zap.L().Named("evse_manager"),
	}
}

func getKey(evseId int) string {
	return fmt.Sprintf("evse-%d", evseId)
}

func (m *managerImpl) InitAll(ctx context.Context) error {
	m.logger.Info("Initializing EVSEs")

	// Create EVSEs from settings stored in the database.
	for _, c := range database.GetEvseSettings(m.db) {
		addErr := m.addEVSEFromSettings(ctx, c)
		if addErr != nil {
			return addErr
		}
	}

	return nil
}

func (m *managerImpl) GetEVSEs() []evse.EVSE {
	m.logger.Debug("Getting all EVSEs")
	var connectors []evse.EVSE

	m.connectors.Range(func(key, value interface{}) bool {
		c, canCast := value.(evse.EVSE)
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

func (m *managerImpl) GetEVSE(evseId int) (evse.EVSE, error) {
	m.logger.Debug("Getting EVSE", zap.Int("evseId", evseId))

	c, isFound := m.connectors.Load(getKey(evseId))
	if isFound {
		return c.(evse.EVSE), nil
	}

	return nil, ErrConnectorNotFound
}

func (m *managerImpl) GetAvailableEVSE() (evse.EVSE, error) {
	m.logger.Debug("Getting available EVSEs")
	var availableConnector evse.EVSE

	m.connectors.Range(func(key, value interface{}) bool {
		c, canCast := value.(evse.EVSE)
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

func (m *managerImpl) StartCharging(evseId int, connectorId *int, measurands []types.Measurand, sampleInterval string) error {
	m.logger.Debug("Attempting to start charging",
		zap.Int("evseId", evseId),
		zap.Any("connectorId", connectorId),
		zap.Any("measurands", measurands),
		zap.String("sampleInterval", sampleInterval),
	)

	c, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	return c.StartCharging(connectorId, measurands, sampleInterval)
}

func (m *managerImpl) StopCharging(evseId int, connectorId *int, reason core.Reason) error {
	m.logger.Debug("Attempting to stop charging",
		zap.Int("evseId", evseId),
		zap.Any("connectorId", connectorId),
		zap.String("reason", string(reason)),
	)

	c, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	return c.StopCharging(reason)
}

func (m *managerImpl) StopAllEVSEs(reason core.Reason) error {
	m.logger.Debug("Stopping all evses", zap.String("reason", string(reason)))

	var err error

	for _, c := range m.GetEVSEs() {
		stopErr := c.StopCharging(reason)
		if stopErr != nil {
			err = stopErr
		}
	}

	return err
}

func (m *managerImpl) AddEVSE(ctx context.Context, c evse.EVSE) error {
	if util.IsNilInterfaceOrPointer(c) {
		return ErrConnectorNil
	}

	err := c.Init(ctx)
	if err != nil {
		return err
	}

	m.logger.Debug("Adding an EVSE to manager", zap.Int("evseId", c.GetEvseId()))

	c.SetNotificationChannel(m.notificationChannel)
	c.SetMeterValuesChannel(m.meterValuesChannel)

	// Add the connector
	m.connectors.Store(getKey(c.GetEvseId()), c)
	return nil
}

func (m *managerImpl) addEVSEFromSettings(ctx context.Context, c settings.EVSE) error {
	logInfo := m.logger.With(zap.Int("evseId", c.EvseId))
	logInfo.Debug("Creating an Evcc from settings")

	// Create EVSE from settings
	evccFromType, err := evcc.NewEVCCFromType(c.EVCC)
	if err != nil {
		return fmt.Errorf("cannot create evcc from type: %w", err)
	}

	// Create a PowerMeter from settings
	logInfo.Debug("Creating power meter")
	meter, powerMeterErr := powerMeter.NewPowerMeter(c.PowerMeter)
	switch {
	case powerMeterErr == nil:
	case errors.Is(powerMeterErr, powerMeter.ErrPowerMeterDisabled):
		logInfo.With(zap.Error(powerMeterErr)).Warn("Power meter disabled")
	case errors.Is(powerMeterErr, powerMeter.ErrPowerMeterUnsupported), errors.Is(powerMeterErr, powerMeter.ErrInvalidConnectionSettings):
		fallthrough
	default:
		logInfo.With(zap.Error(powerMeterErr)).Error("Cannot instantiate power meter for evse")
		return err
	}

	// Create EVSE from EVCC and Power Meter
	logInfo.Debug("Creating EVSE")
	evse, err := evse.NewEvse(zap.L(), c.EvseId, evccFromType, meter, float64(c.MaxPower), nil)
	if err != nil {
		return err
	}

	if m.notificationChannel != nil {
		evse.SetNotificationChannel(m.notificationChannel)
	}

	return m.AddEVSE(ctx, evse)
}

func (m *managerImpl) UpdateEVSE(ctx context.Context, c evse.EVSE) error {
	m.logger.With(zap.Int("evseId", c.GetEvseId())).Debug("Updating an EVSE")
	// todo implement me
	return nil
}

func (m *managerImpl) RemoveEVSE(evseId int) error {
	m.logger.With(zap.Int("evseId", evseId)).Debug("Removing an EVSE")

	m.connectors.Delete(getKey(evseId))
	return nil
}

func (m *managerImpl) RestoreEVSEs() error {
	m.logger.Debug("Attempting to restore EVSEs")

	for _, s := range database.GetEvseSettings(m.db) {
		err := m.restoreEVSEStatus(s)
		if err != nil {
			m.logger.With(zap.Error(err), zap.Int("evse_id", s.EvseId)).Error("Error restoring an EVSE")
			continue
		}
	}

	return nil
}

func (m *managerImpl) restoreEVSEStatus(c settings.EVSE) error {
	logInfo := m.logger.With(zap.Any("settings", c))
	logInfo.Debug("Attempting to restore connector status")

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

		logInfo.Debug("Successfully resumed charging")
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

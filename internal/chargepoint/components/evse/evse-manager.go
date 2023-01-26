package evse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/session"
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
		FindEVSE(evseId int) (EVSE, error)
		FindAvailableEVSE() (EVSE, error)
		FindEVSEWithTagId(tagId string) (EVSE, error)
		FindEVSEWithTransactionId(transactionId string) (EVSE, error)
		FindEVSEWithReservationId(reservationId int) (EVSE, error)

		StartCharging(evseId int, tagId, transactionId string) error
		StopCharging(tagId, transactionId string, reason core.Reason) error
		StopAllEVSEs(reason core.Reason) error

		ImportFromSettings(c []settings.EVSE) error
		RestoreEVSEStatus(*settings.EVSE) error

		SetMaxChargingTime(maxChargingTime int) error
		SetNotificationChannel(notificationChannel chan notifications.StatusNotification)
		GetNotificationChannel() chan notifications.StatusNotification
		SetMeterValuesChannel(notificationChannel chan notifications.MeterValueNotification)
	}

	managerImpl struct {
		// Used to store settings of EVSEs
		db *badger.DB
		// Used for running EVSE instances
		connectors          sync.Map
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
	var evseSettings []settings.EVSE

	// Query the database for EVSE settings.
	err := m.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte("evse-")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var data settings.EVSE
			item := it.Item()

			// Value should be the EVSE struct.
			err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &data)
			})
			if err != nil {
				continue
			}
		}
		return txn.Commit()
	})
	if err != nil {
		log.WithError(err).Error("Error querying for EVSE settings")
		return err
	}

	// Create EVSEs from settings stored in the database.
	for _, c := range evseSettings {
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

func (m *managerImpl) FindEVSE(evseId int) (EVSE, error) {
	c, isFound := m.connectors.Load(getKey(evseId))
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

func (m *managerImpl) AddEVSE(ctx context.Context, c EVSE) error {
	if util.IsNilInterfaceOrPointer(c) {
		return ErrConnectorNil
	}

	// todo retry mechanism?
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
	evse, err := NewEvse(c.EvseId, evccFromType, meter, c.PowerMeter.Enabled, float64(c.MaxPower), nil)
	if err != nil {
		return err
	}

	return m.AddEVSE(ctx, evse)
}

func (m *managerImpl) ImportFromSettings(connectors []settings.EVSE) error {
	log.Info("Importing connectors to the database")
	// Sync the settings to the database
	return m.db.Update(func(txn *badger.Txn) error {

		for _, connector := range connectors {
			marshal, err := json.Marshal(connector)
			if err != nil {
				return err
			}

			err = txn.Set([]byte(getKey(connector.EvseId)), marshal)
			if err != nil {
				return err
			}
		}

		return txn.Commit()
	})
}

func (m *managerImpl) UpdateEVSE(ctx context.Context, c EVSE) error {
	return nil
}

func (m *managerImpl) RemoveEVSE(evseId int) error {
	return nil
}

func (m *managerImpl) SetMaxChargingTime(maxChargingTime int) error {
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
		// todo (c.Session)
		timeLeft, err := evse.ResumeCharging(session.Session{})
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
		_, schedulerErr := scheduler.NewScheduler().Every(*timeLeft).Minutes().At(1).Do(evse.StopCharging(core.ReasonLocal))
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

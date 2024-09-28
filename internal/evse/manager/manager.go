package manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/ChargePi/ChargePi-go/internal/pkg/scheduler"
	"github.com/ChargePi/ChargePi-go/internal/sessions"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"go.uber.org/zap"
)

// Manager is responsible for managing EVSE instances at runtime. It persists their status and settings in the database.
// It also provides methods to start and stop charging, get current consumption and more.

type Impl struct {
	// Used to store settings of EVSEs
	repository     sessions.SessionRepository
	evseRepository EvseSettingsRepository
	// Used for running EVSE instances
	evses               sync.Map
	reservations        *reservations
	notificationChannel chan notifications.StatusNotification
	meterValuesChannel  chan notifications.MeterValueNotification
	logger              *zap.Logger
}

func NewManager(sessionRepository sessions.SessionRepository, settingsRepository EvseSettingsRepository, notificationChannel chan notifications.StatusNotification) (*Impl, error) {
	return &Impl{
		repository:          sessionRepository,
		evseRepository:      settingsRepository,
		evses:               sync.Map{},
		reservations:        newReservationsStore(),
		notificationChannel: notificationChannel,
		logger:              zap.L().Named("evse_manager"),
	}, nil
}

func getKey(evseId int) string {
	return fmt.Sprintf("evse-%d", evseId)
}

// InitAll fetches all EVSE settings from the database, creates EVSE instances and initializes them.
// This method should be called at startup.
func (m *Impl) InitAll(ctx context.Context) error {
	m.logger.Info("Initializing EVSEs")

	settings, err := m.evseRepository.GetEvseSettings()
	if err != nil {
		return err
	}

	// Create EVSEs from settings
	for _, c := range settings {
		addErr := m.AddEVSEFromSettings(ctx, c)
		if addErr != nil {
			return addErr
		}
	}

	return nil
}

// AddEVSE adds an EVSE instance to the manager. This call should only be called internally, not exposed to APIs.
func (m *Impl) AddEVSE(ctx context.Context, evse evse.EVSE) error {
	if util.IsNilInterfaceOrPointer(evse) {
		return ErrConnectorNil
	}

	// todo check if evse already initialized - skip the step if so
	// Initialize the EVSE
	err := evse.Init(ctx)
	if err != nil {
		return err
	}

	m.logger.Debug("Adding an EVSE",
		zap.Int("evseId", evse.GetEvseId()))

	// Override the notification channels
	evse.SetNotificationChannel(m.notificationChannel)
	evse.SetMeterValuesChannel(m.meterValuesChannel)

	// Add the connector
	m.evses.Store(getKey(evse.GetEvseId()), evse)

	// todo Store the settings in the database

	return nil
}

func (m *Impl) AddEVSEFromSettings(ctx context.Context, settings evse.Settings) error {
	m.logger.Debug("Adding an EVSE to manager", zap.Int("evseId", settings.EvseId))

	evse, err := evse.NewEvseFromSettings(settings)
	if err != nil {
		return err
	}

	if m.notificationChannel != nil {
		evse.SetNotificationChannel(m.notificationChannel)
	}

	return m.AddEVSE(ctx, evse)
}

func (m *Impl) GetEVSEs() []evse.EVSE {
	m.logger.Debug("Getting all EVSEs")
	var connectors []evse.EVSE

	m.evses.Range(func(key, value interface{}) bool {
		c, canCast := value.(evse.EVSE)
		if canCast {
			connectors = append(connectors, c)
		}
		return true
	})

	return connectors
}

func (m *Impl) GetNotificationChannel() chan notifications.StatusNotification {
	return m.notificationChannel
}

func (m *Impl) GetEVSE(evseId int) (evse.EVSE, error) {
	m.logger.Debug("Getting EVSE", zap.Int("evseId", evseId))

	c, isFound := m.evses.Load(getKey(evseId))
	if isFound {
		return c.(evse.EVSE), nil
	}

	return nil, ErrConnectorNotFound
}

func (m *Impl) GetAvailableEVSE() (evse.EVSE, error) {
	m.logger.Debug("Getting available EVSEs")
	var availableConnector evse.EVSE

	m.evses.Range(func(key, value interface{}) bool {
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

func (m *Impl) UpdateEVSE(ctx context.Context, c evse.EVSE) error {
	m.logger.Debug("Updating an EVSE",
		zap.Int("evseId", c.GetEvseId()))

	// Store current state in the database, just in case we need to restore it
	// todo implement me
	return nil
}

func (m *Impl) RemoveEVSE(evseId int) error {
	m.logger.Debug("Removing an EVSE",
		zap.Int("evseId", evseId))

	evse, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	// todo Remove all reservations for this EVSE
	// m.reservations.RemoveReservationsForEVSE(evseId)

	// todo Store current state in the database, just in case we need to restore it

	// Shutdown procedure: stop charging, terminate connections and remove the EVSE
	err = evse.Cleanup()
	if err != nil {
		return err
	}

	m.evses.Delete(getKey(evseId))

	// todo remove settings from the database?
	// todo schedule removal of restored EVSE state?
	// todo schedule removal of sessions associated with this EVSE

	return nil
}

func (m *Impl) StartCharging(evseId int, connectorId *int, measurands []types.Measurand, sampleInterval string) error {
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

	// Schedule a stop charging after the maxChargingTime, if provided
	/*if m.maxChargingTime != nil {
		// Schedule a stop charging after the maxChargingTime
		_, err := m.scheduler.Every(*evse.maxChargingTime).Minutes().Tag(fmt.Sprintf("evse-%d-chargingTimer",evseId)).Do(evse.StopCharging, core.ReasonLocal)
		if err != nil {
			logInfo.WithError(err).Error("Cannot schedule stop charging")
		}
	}*/

	return c.StartCharging(connectorId, measurands, sampleInterval)
}

func (m *Impl) StopCharging(evseId int, connectorId *int, reason core.Reason) error {
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

func (m *Impl) StopAll(reason core.Reason) error {
	m.logger.Debug("Stopping all evses", zap.String("reason", string(reason)))

	var err error

	// todo Store each state in the database, just in case we need to restore it
	for _, c := range m.GetEVSEs() {
		stopErr := c.StopCharging(reason)
		if stopErr != nil {
			err = stopErr
		}
	}

	return err
}

func (m *Impl) GetCurrentConsumption(evseId, connectorId *int) (*types.MeterValue, error) {
	// Get consumption on all connectors
	if evseId == nil {
		return nil, nil
	}

	// Get consumption on a specific EVSE
	evse, err := m.GetEVSE(*evseId)
	if err != nil {
		return nil, err
	}

	energy, err := evse.GetPowerMeter().GetEnergy()
	if err != nil {
		return nil, err
	}

	return &types.MeterValue{
		Timestamp: types.Now(),
		SampledValue: []types.SampledValue{
			*energy,
		},
	}, nil
}

func (m *Impl) UnlockConnector(evseId, connectorId int) error {
	m.logger.Debug("Unlocking connector",
		zap.Int("evseId", evseId),
		zap.Any("connectorId", connectorId),
	)

	evse, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	evse.Unlock(connectorId)
	return nil
}

func (m *Impl) RestoreEVSEs() error {
	m.logger.Debug("Attempting to restore EVSEs")

	// Get all EVSE settings from the database
	settings, err := m.evseRepository.GetEvseSettings()
	if err != nil {
		return err
	}

	errChan := make(chan error, len(settings))

	var wg sync.WaitGroup
	wg.Add(len(settings))

	// Restore each EVSE concurrently
	for _, s := range settings {
		go func() {
			defer wg.Done()
			err := m.restoreEVSEStatus(s)
			if err != nil {
				m.logger.With(zap.Error(err), zap.Int("evse_id", s.EvseId)).Error("Error restoring an EVSE")
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	// todo Aggregate the errors?
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Impl) restoreEVSEStatus(c evse.Settings) error {
	logInfo := m.logger.With(zap.Any("settings", c))
	logInfo.Debug("Attempting to restore EVSE from status")

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

func (m *Impl) GetEVSEWithReservationId(reservationId int) (evse.EVSE, error) {
	m.logger.With(zap.Int("reservationId", reservationId)).Debug("Finding evse with reservation id")

	// Get reservation details
	reservationDetails, err := m.reservations.GetEVSEWithReservationId(reservationId)
	if err != nil {
		return nil, err
	}

	// Get evse from reservation details
	evse, err := m.GetEVSE(reservationDetails.EvseId)
	if err != nil {
		return nil, err
	}

	return evse, nil
}

func (m *Impl) Reserve(evseId int, connectorId *int, reservationId int, tagId string) error {
	m.logger.Debug("Reserving EVSE",
		zap.Int("evseId", evseId),
		zap.Intp("connectorId", connectorId),
		zap.Int("reservationId", reservationId),
		zap.String("tagId", tagId),
	)

	// Verify that the evse exists
	evse, err := m.GetEVSE(evseId)
	if err != nil {
		return err
	}

	details := reservation{
		EvseId:      evseId,
		ConnectorId: connectorId,
		TagId:       tagId,
	}
	err = m.reservations.Reserve(reservationId, details)
	if err != nil {
		return err
	}

	// Set evse status to reserved
	err = evse.SetStatus(core.ChargePointStatusReserved, core.NoError)
	if err != nil {
		return err
	}

	return nil
}

func (m *Impl) RemoveReservation(reservationId int) error {
	m.logger.With(zap.Int("reservationId", reservationId)).Debug("Removing reservation")

	// Get EVSE from reservation
	evse, err := m.GetEVSEWithReservationId(reservationId)
	if err != nil {
		return err
	}

	// Remove reservation
	err = m.reservations.RemoveReservation(reservationId)
	if err != nil {
		return err
	}

	err = evse.SetStatus(core.ChargePointStatusAvailable, core.NoError)
	if err != nil {
		return err
	}

	return nil
}

func (m *Impl) Shutdown() error {
	m.logger.Info("Shutting down manager")
	// Stop charging all
	err := m.StopAll(core.ReasonLocal)
	if err != nil {
		return err
	}

	// Clean up all EVSEs runtime
	for _, e := range m.GetEVSEs() {
		err := e.Cleanup()
		if err != nil {
			return err
		}
	}

	// Close the notification channel
	close(m.notificationChannel)
	close(m.meterValuesChannel)

	return nil
}

package manager

import (
	"context"

	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/pkg/errors"
)

var (
	ErrConnectorNotFound      = errors.New("connector not found")
	ErrReservationNotFound    = errors.New("reservation not found")
	ErrConnectorStatusInvalid = errors.New("connector status invalid")
	ErrConnectorNil           = errors.New("connector is nil")
)

type Manager interface {
	InitAll(ctx context.Context) error
	AddEVSE(ctx context.Context, c evse.EVSE) error
	AddEVSEFromSettings(ctx context.Context, settings evse.Settings) error
	UpdateEVSE(ctx context.Context, c evse.EVSE) error
	RemoveEVSE(evseId int) error
	GetEVSEs() []evse.EVSE
	GetEVSE(evseId int) (evse.EVSE, error)
	GetAvailableEVSE() (evse.EVSE, error)

	StartCharging(evseId int, connectorId *int, measurands []types.Measurand, sampleInterval string) error
	StopCharging(evseId int, connectorId *int, reason core.Reason) error
	StopAll(reason core.Reason) error

	GetCurrentConsumption(evseId, connectorId *int) (*types.MeterValue, error)
	UnlockConnector(evseId, connectorId int) error

	GetEVSEWithReservationId(reservationId int) (evse.EVSE, error)
	Reserve(evseId int, connectorId *int, reservationId int, tagId string) error
	RemoveReservation(reservationId int) error

	GetNotificationChannel() chan notifications.StatusNotification
	// GetMeterValuesChannel() chan notifications.MeterValueNotification

	Shutdown() error
}

type EvseSettingsRepository interface {
	GetEvseSettings() ([]evse.Settings, error)
	SetEvseSettings([]evse.Settings) error
}

package evse

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/pkg/errors"

	"github.com/ChargePi/ChargePi-go/internal/pkg/notifications"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/evcc"
	powerMeter "github.com/ChargePi/ChargePi-go/pkg/hardware/power-meter"
)

var (
	ErrInvalidEvseId        = errors.New("invalid evse id")
	ErrInvalidStatus        = errors.New("invalid evse status")
	ErrInvalidEVCC          = errors.New("evcc cannot be nil")
	ErrNotCharging          = errors.New("evse not charging")
	ErrPowerMeterNotEnabled = errors.New("power meter not enabled")
	ErrConnectorExists      = errors.New("connector already exists")
)

type EVSE interface {
	Init(ctx context.Context) error
	Cleanup() error

	AddConnector(connector ConnectorSettings) error
	GetConnectors() []ConnectorSettings

	GetMaxChargingPower() float64
	GetEvseId() int

	StartCharging(connectorId *int, measurands []types.Measurand, sampleInterval string) error
	StopCharging(reason core.Reason) error
	Lock(connectorId int)
	Unlock(connectorId int)

	GetPowerMeter() powerMeter.PowerMeter
	SetPowerMeter(powerMeter.PowerMeter) error
	SamplePowerMeter(measurands []types.Measurand) ([]types.SampledValue, error)

	GetEvcc() evcc.EVCC
	SetEvcc(evcc.EVCC) error

	SetAvailability(isAvailable bool) error
	SetStatus(status core.ChargePointStatus, errCode core.ChargePointErrorCode) error
	GetStatus() (core.ChargePointStatus, core.ChargePointErrorCode)

	IsAvailable() bool
	IsPreparing() bool
	IsCharging() bool
	IsReserved() bool
	IsUnavailable() bool

	SetNotificationChannel(notificationChannel chan<- notifications.StatusNotification)
	SetMeterValuesChannel(notificationChannel chan<- notifications.MeterValueNotification)
	IsHealthy() bool
}

type Info struct {
	EvseId   int
	MaxPower float64
}

type ConnectorSettings struct {
	ConnectorId int    `json:"connectorId,omitempty" yaml:"connectorId" mapstructure:"connectorId" validate:"required,gte=0"`
	Type        string `json:"type,omitempty" yaml:"type" mapstructure:"type" validate:"required,oneof=Type1 Type2 CCS CHAdeMO Tesla CCS2"`
	Status      string `json:"status,omitempty" yaml:"status" mapstructure:"status"`
}

func (c *ConnectorSettings) Validate() error {
	return validator.New().Struct(c)
}

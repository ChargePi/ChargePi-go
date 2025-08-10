package powerMeter

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

// Supported power meters
const (
	TypeC5460A = "cs5460a"
	TypeDummy  = "dummy"
)

var (
	ErrPowerMeterUnsupported     = errors.New("power meter type not supported")
	ErrPowerMeterDisabled        = errors.New("power meter not enabled")
	ErrInvalidConnectionSettings = errors.New("invalid power meter connection settings")
)

// PowerMeter is an abstraction for measurement hardware.
type PowerMeter interface {
	Init(ctx context.Context) error
	Cleanup()
	Reset()
	GetEnergy() (*types.SampledValue, error)
	GetPower(phase int) (*types.SampledValue, error)
	GetReactivePower(phase int) (*types.SampledValue, error)
	GetApparentPower(phase int) (*types.SampledValue, error)
	GetCurrent(phase int) (*types.SampledValue, error)
	GetVoltage(phase int) (*types.SampledValue, error)
	GetType() string
}

// NewPowerMeter creates a new power meter based on the connector settings.
func NewPowerMeter(meterSettings settings.PowerMeter) (PowerMeter, error) {
	logger := zap.L()
	if meterSettings.Enabled {
		logger.With(zap.String("meter_type", meterSettings.Type)).Info("Creating a new power meter")

		switch meterSettings.Type {
		case TypeC5460A:
			if util.IsNilInterfaceOrPointer(meterSettings.SPI) {
				return nil, ErrInvalidConnectionSettings
			}

			powerMeter, err := NewCS5460PowerMeter(
				logger,
				meterSettings.SPI.ChipSelect,
				meterSettings.SPI.Bus,
				meterSettings.CS5460.ShuntOffset,
				meterSettings.CS5460.VoltageDividerOffset,
			)
			if err != nil {
				return nil, err
			}

			return powerMeter, nil
		case TypeDummy:
			return NewDummy(logger, meterSettings.PowerMeterDummy)
		default:
			return nil, ErrPowerMeterUnsupported
		}
	}

	return nil, ErrPowerMeterDisabled
}

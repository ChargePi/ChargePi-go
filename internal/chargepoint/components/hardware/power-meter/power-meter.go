package powerMeter

import (
	"context"
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

// Supported power meters
const (
	TypeC5460A = "cs5460a"
)

var (
	ErrPowerMeterUnsupported     = errors.New("power meter type not supported")
	ErrPowerMeterDisabled        = errors.New("power meter not enabled")
	ErrInvalidConnectionSettings = errors.New("invalid power meter connection settings")
)

// PowerMeter is an abstraction for measurement hardware.
type (
	PowerMeter interface {
		Init(ctx context.Context) error
		Reset()
		GetEnergy() types.SampledValue
		GetPower() types.SampledValue
		GetCurrent(phase int) types.SampledValue
		GetVoltage(phase int) types.SampledValue
		GetRMSCurrent(phase int) types.SampledValue
		GetRMSVoltage(phase int) types.SampledValue
		GetType() string
	}
)

// NewPowerMeter creates a new power meter based on the connector settings.
func NewPowerMeter(meterSettings settings.PowerMeter) (PowerMeter, error) {
	if meterSettings.Enabled {
		log.Infof("Creating a new power meter: %s", meterSettings.Type)

		switch meterSettings.Type {
		case TypeC5460A:
			if util.IsNilInterfaceOrPointer(meterSettings.SPI) {
				return nil, ErrInvalidConnectionSettings
			}

			powerMeter, err := NewCS5460PowerMeter(
				meterSettings.SPI.ChipSelect,
				meterSettings.SPI.Bus,
				meterSettings.CS5460.ShuntOffset,
				meterSettings.CS5460.VoltageDividerOffset,
			)
			if err != nil {
				return nil, err
			}

			return powerMeter, nil
		default:
			return nil, ErrPowerMeterUnsupported
		}
	}

	return nil, ErrPowerMeterDisabled
}

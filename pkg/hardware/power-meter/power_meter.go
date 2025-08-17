package powerMeter

import (
	"context"
	"errors"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

var (
	ErrPowerMeterUnsupported     = errors.New("power meter type not supported")
	ErrPowerMeterDisabled        = errors.New("power meter not enabled")
	ErrInvalidConnectionSettings = errors.New("invalid power meter connection settings")
)

// PowerMeter is an abstraction for measurement hardware.
type PowerMeter interface {
	Init(ctx context.Context) error
	Cleanup() error
	Reset()
	GetEnergy() (*types.SampledValue, error)
	GetPower(phase int) (*types.SampledValue, error)
	GetReactivePower(phase int) (*types.SampledValue, error)
	GetApparentPower(phase int) (*types.SampledValue, error)
	GetCurrent(phase int) (*types.SampledValue, error)
	GetVoltage(phase int) (*types.SampledValue, error)
	GetType() string
}

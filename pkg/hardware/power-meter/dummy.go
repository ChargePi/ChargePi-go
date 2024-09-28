package powerMeter

import (
	"context"
	"go.uber.org/zap"
	"sync/atomic"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

const TypeDummy = "dummy"

type DummySettings struct {
	Voltage          float64 `json:"voltage,omitempty" yaml:"voltage,omitempty" mapstructure:"voltage,omitempty"`
	VoltageBehaviour string  `json:"voltageBehaviour,omitempty" yaml:"voltageBehaviour,omitempty" mapstructure:"voltageBehaviour,omitempty"`
	BaseCurrent      float64 `json:"baseCurrent,omitempty" yaml:"baseCurrent,omitempty" mapstructure:"baseCurrent,omitempty"`
	CurrentBehaviour string  `json:"currentBehaviour,omitempty" yaml:"currentBehaviour,omitempty" mapstructure:"currentBehaviour,omitempty"`
}

type Dummy struct {
	logger     *zap.Logger
	settings   DummySettings
	energy     atomic.Int64
	lastSample *time.Time
}

func NewDummy(logger *zap.Logger, settings *DummySettings) (*Dummy, error) {
	if settings == nil {
		return nil, ErrInvalidConnectionSettings
	}

	return &Dummy{
		settings: *settings,
		logger:   logger.Named("power_meter_dummy"),
		energy:   atomic.Int64{},
	}, nil
}

func (d *Dummy) Init(ctx context.Context) error {
	d.logger.Info("Initializing power meter")
	return nil
}

func (d *Dummy) Reset() {
	d.logger.Info("Resetting power meter")
}

func (d *Dummy) GetEnergy() (*types.SampledValue, error) {
	d.logger.Info("Getting energy")
	return nil, nil
}

func (d *Dummy) GetPower(phase int) (*types.SampledValue, error) {
	d.logger.With(zap.Int("phase", phase)).Info("Getting power")
	return nil, nil
}

func (d *Dummy) GetReactivePower(phase int) (*types.SampledValue, error) {
	d.logger.With(zap.Int("phase", phase)).Info("Getting reactive power")
	return nil, nil
}

func (d *Dummy) GetApparentPower(phase int) (*types.SampledValue, error) {
	d.logger.With(zap.Int("phase", phase)).Info("Getting apparent power")
	return nil, nil
}

func (d *Dummy) GetCurrent(phase int) (*types.SampledValue, error) {
	d.logger.With(zap.Int("phase", phase)).Info("Getting current")
	return nil, nil
}

func (d *Dummy) GetVoltage(phase int) (*types.SampledValue, error) {
	d.logger.With(zap.Int("phase", phase)).Info("Getting voltage")
	return nil, nil
}

func (d *Dummy) GetType() string {
	return TypeDummy
}

func (d *Dummy) Cleanup() error {
	d.logger.Info("Cleaning up power meter")
	return nil
}

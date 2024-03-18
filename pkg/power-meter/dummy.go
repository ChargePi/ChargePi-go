package powerMeter

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type Dummy struct {
	logger     log.FieldLogger
	settings   settings.PowerMeterDummy
	energy     atomic.Int64
	lastSample *time.Time
}

func NewDummy(settings *settings.PowerMeterDummy) (*Dummy, error) {
	if settings == nil {
		return nil, ErrInvalidConnectionSettings
	}

	return &Dummy{
		settings: *settings,
		logger:   log.StandardLogger().WithField("component", "power-meter-dummy"),
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
	d.logger.WithField("phase", phase).Info("Getting power")

	return nil, nil
}

func (d *Dummy) GetReactivePower(phase int) (*types.SampledValue, error) {
	d.logger.WithField("phase", phase).Info("Getting reactive power")
	return nil, nil
}

func (d *Dummy) GetApparentPower(phase int) (*types.SampledValue, error) {
	d.logger.WithField("phase", phase).Info("Getting apparent power")
	return nil, nil
}

func (d *Dummy) GetCurrent(phase int) (*types.SampledValue, error) {
	d.logger.WithField("phase", phase).Info("Getting current")

	return nil, nil
}

func (d *Dummy) GetVoltage(phase int) (*types.SampledValue, error) {
	d.logger.WithField("phase", phase).Info("Getting voltage")
	return nil, nil
}

func (d *Dummy) GetType() string {
	return TypeDummy
}

func (d *Dummy) Cleanup() {
	d.logger.Info("Cleaning up power meter")
}

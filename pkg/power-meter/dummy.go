package powerMeter

import (
	"context"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type Dummy struct {
	logger     *log.Logger
	settings   settings.PowerMeterDummy
	energy     int64
	lastSample *time.Time
}

func NewDummy(settings *settings.PowerMeterDummy) (*Dummy, error) {
	return &Dummy{
		settings: *settings,
		logger:   log.StandardLogger(),
		energy:   0,
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

	return nil, nil
}

func (d *Dummy) GetPower(phase int) (*types.SampledValue, error) {
	return nil, nil
}

func (d *Dummy) GetReactivePower(phase int) (*types.SampledValue, error) {
	return nil, nil
}

func (d *Dummy) GetApparentPower(phase int) (*types.SampledValue, error) {
	return nil, nil
}

func (d *Dummy) GetCurrent(phase int) (*types.SampledValue, error) {
	return nil, nil
}

func (d *Dummy) GetVoltage(phase int) (*types.SampledValue, error) {
	return nil, nil
}

func (d *Dummy) GetType() string {
	return TypeDummy
}

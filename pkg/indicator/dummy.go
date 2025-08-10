package indicator

import (
	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"go.uber.org/zap"
)

type Dummy struct {
	logger     *zap.Logger
	settings   settings.IndicatorDummy
	brightness int
}

func NewDummy(logger *zap.Logger, settings *settings.IndicatorDummy) *Dummy {
	return &Dummy{
		settings: *settings,
		logger:   logger.Named("indicator_dummy"),
	}
}

func (d *Dummy) ChangeColor(index int, color Color) error {
	d.logger.Sugar().Infof("Changing color of index %d to %s", index, color)
	return nil
}

func (d *Dummy) Blink(index int, times int, color Color) error {
	d.logger.Sugar().Infof("Blinking color of index %d %d times to %s", index, times, color)
	return nil
}

func (d *Dummy) SetBrightness(brightness int) error {
	d.logger.Sugar().Infof("Setting brightness to %d", brightness)
	d.brightness = brightness
	return nil
}

func (d *Dummy) GetBrightness() int {
	d.logger.Info("Getting brightness")
	return d.brightness
}

func (d *Dummy) Cleanup() {
	d.logger.Info("Cleaning up indicator")
}

func (d *Dummy) GetType() string {
	return TypeDummy
}

package display

import (
	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"go.uber.org/zap"
)

type Dummy struct {
	logger   *zap.Logger
	settings settings.DisplayDummy
}

func NewDummy(logger *zap.Logger, settings *settings.DisplayDummy) (*Dummy, error) {
	return &Dummy{
		logger:   logger.Named("display_dummy"),
		settings: *settings,
	}, nil
}

func (d *Dummy) DisplayMessage(message display.MessageInfo) {
	d.logger.With(zap.Any("message", message)).Info("Displaying message")
}

func (d *Dummy) Clear() {
	d.logger.Info("Clearing display")
}

func (d *Dummy) Cleanup() {
	d.logger.Info("Cleaning up display")
}

func (d *Dummy) GetType() string {
	return TypeDummy
}

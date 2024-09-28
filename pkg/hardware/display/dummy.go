package display

import (
	"context"

	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"go.uber.org/zap"
)

const TypeDummy = "dummy"

type DummySettings struct{}

type Dummy struct {
	logger   *zap.Logger
	settings DummySettings
}

func NewDummy(logger *zap.Logger, settings *DummySettings) (*Dummy, error) {
	return &Dummy{
		logger:   logger.Named("display_dummy"),
		settings: *settings,
	}, nil
}

func (d *Dummy) DisplayMessage(message display.MessageInfo) error {
	d.logger.With(zap.Any("message", message)).Info("Displaying message")
	return nil
}

func (d *Dummy) Clear() error {
	d.logger.Info("Clearing display")
	return nil
}

func (d *Dummy) Cleanup(context.Context) error {
	d.logger.Info("Cleaning up display")
	return nil
}

func (d *Dummy) GetType() string {
	return TypeDummy
}

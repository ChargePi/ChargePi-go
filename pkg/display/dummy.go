package display

import (
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type Dummy struct {
	logger   *log.Logger
	settings settings.DisplayDummy
}

func NewDummy(settings *settings.DisplayDummy) (*Dummy, error) {
	return &Dummy{
		settings: *settings,
	}, nil
}

func (d *Dummy) DisplayMessage(message display.MessageInfo) {
	d.logger.WithFields(log.Fields{"message": message}).Info("Displaying message")
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

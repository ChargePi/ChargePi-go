package hardware

import (
	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
)

type (
	RelayImpl struct {
		RelayPin     int
		InverseLogic bool
		currentState bool
		pin          *gpiod.Line
	}

	Relay interface {
		Enable()
		Disable()
	}
)

// NewRelay creates a new RelayImpl struct that will communicate with the GPIO pin specified.
func NewRelay(relayPin int, inverseLogic bool) *RelayImpl {
	if relayPin <= 0 {
		return nil
	}

	log.Debugf("Creating new relay at pin %d", relayPin)
	relay := RelayImpl{
		RelayPin:     relayPin,
		InverseLogic: inverseLogic,
		currentState: inverseLogic,
	}

	err := relay.initPin()
	if err != nil {
		return nil
	}

	return &relay
}

func (r *RelayImpl) initPin() error {
	// Refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}

	r.pin, err = c.RequestLine(r.RelayPin, gpiod.AsOutput(0))
	return err
}

func (r *RelayImpl) Enable() {
	if r.InverseLogic {
		_ = r.pin.SetValue(0)
	} else {
		_ = r.pin.SetValue(1)
	}

	// Always consider positive logic for status determination
	r.currentState = true
}

func (r *RelayImpl) Disable() {
	if r.InverseLogic {
		_ = r.pin.SetValue(1)
	} else {
		_ = r.pin.SetValue(0)
	}

	// Always consider positive logic for status determination
	r.currentState = false
}

func (r *RelayImpl) Close() {
	_ = r.pin.Close()
}

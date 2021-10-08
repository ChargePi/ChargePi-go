package hardware

import (
	"github.com/warthog618/gpiod"
)

type Relay struct {
	RelayPin     int
	InverseLogic bool
	currentState bool
	pin          *gpiod.Line
}

// NewRelay creates a new Relay struct that will communicate with the GPIO pin specified.
func NewRelay(relayPin int, inverseLogic bool) *Relay {
	if relayPin <= 0 {
		return nil
	}
	var relay = Relay{
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

func (receiver *Relay) initPin() error {
	//refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}
	receiver.pin, err = c.RequestLine(receiver.RelayPin, gpiod.AsOutput(0))
	if err != nil {
		return err
	}
	return nil
}

func (receiver *Relay) On() {
	if receiver.InverseLogic {
		receiver.pin.SetValue(0)
	} else {
		receiver.pin.SetValue(1)
	}
	// Always consider positive logic for status determination
	receiver.currentState = true
}

func (receiver *Relay) Off() {
	if receiver.InverseLogic {
		receiver.pin.SetValue(1)
	} else {
		receiver.pin.SetValue(0)
	}
	// Always consider positive logic for status determination
	receiver.currentState = false
}

func (receiver *Relay) Close() {
	receiver.pin.Close()
}

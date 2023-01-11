package evcc

import (
	"context"
	"errors"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
)

var (
	ErrInvalidPinNumber = errors.New("pin number must be greater than 0")
)

type RelayAsEvcc struct {
	relayPin      int
	inverseLogic  bool
	state         CarState
	pin           *gpiod.Line
	statusChannel chan StateNotification
}

// NewRelay creates a new RelayImpl struct that will communicate with the GPIO pin specified.
func NewRelay(relayPin int, inverseLogic bool) (*RelayAsEvcc, error) {
	if relayPin <= 0 {
		return nil, ErrInvalidPinNumber
	}

	log.Debugf("Creating new relay at pin %d", relayPin)
	relay := RelayAsEvcc{
		relayPin:      relayPin,
		inverseLogic:  inverseLogic,
		statusChannel: make(chan StateNotification, 10),
	}

	return &relay, nil
}

func (r *RelayAsEvcc) Init(ctx context.Context) error {
	// Refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}

	r.pin, err = c.RequestLine(r.relayPin, gpiod.AsOutput(0))
	return err
}

func (r *RelayAsEvcc) Lock() {
}

func (r *RelayAsEvcc) Unlock() {
}

func (r *RelayAsEvcc) GetError() string {
	return string(core.NoError)
}

func (r *RelayAsEvcc) GetState() CarState {
	return r.state
}

func (r *RelayAsEvcc) EnableCharging() error {
	if r.inverseLogic {
		_ = r.pin.SetValue(0)
	} else {
		_ = r.pin.SetValue(1)
	}

	_ = r.setState(StateB2, string(core.NoError))
	return r.setState(StateC2, string(core.NoError))
}

func (r *RelayAsEvcc) DisableCharging() {
	if r.inverseLogic {
		_ = r.pin.SetValue(1)
	} else {
		_ = r.pin.SetValue(0)
	}

	_ = r.setState(StateB1, string(core.NoError))
	_ = r.setState(StateA1, string(core.NoError))
}

func (r *RelayAsEvcc) setState(state CarState, error string) error {
	if !IsStateValid(state) {
		return nil
	}

	r.state = state
	r.statusChannel <- NewStateNotification(state, error)
	return nil
}

func (r *RelayAsEvcc) SetMaxChargingCurrent(value float64) error {
	return nil
}

func (r *RelayAsEvcc) GetMaxChargingCurrent() float64 {
	return 0.0
}

func (r *RelayAsEvcc) Cleanup() error {
	return r.pin.Close()
}

func (r *RelayAsEvcc) GetType() string {
	return Relay
}

func (r *RelayAsEvcc) GetStatusChangeChannel() <-chan StateNotification {
	return r.statusChannel
}

func (r *RelayAsEvcc) Reset() {
}

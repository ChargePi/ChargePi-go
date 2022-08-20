package evcc

import (
	"context"
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/evcc"
)

var (
	ErrInvalidPinNumber = errors.New("pin number must be greater than 0")
	ErrInitFailed       = errors.New("init failed")
)

type (
	relayAsEvcc struct {
		relayPin      int
		inverseLogic  bool
		state         evcc.CarState
		pin           *gpiod.Line
		statusChannel chan evcc.StateNotification
	}
)

// NewRelay creates a new RelayImpl struct that will communicate with the GPIO pin specified.
func NewRelay(relayPin int, inverseLogic bool) (*relayAsEvcc, error) {
	if relayPin <= 0 {
		return nil, ErrInvalidPinNumber
	}

	log.Debugf("Creating new relay at pin %d", relayPin)
	relay := relayAsEvcc{
		relayPin:      relayPin,
		inverseLogic:  inverseLogic,
		statusChannel: make(chan evcc.StateNotification, 10),
	}

	return &relay, nil
}

func (r *relayAsEvcc) Init(ctx context.Context) error {
	// Refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}

	r.pin, err = c.RequestLine(r.relayPin, gpiod.AsOutput(0))
	return err
}

func (r *relayAsEvcc) Lock() {
}

func (r *relayAsEvcc) Unlock() {
}

func (r *relayAsEvcc) GetError() string {
	return string(core.NoError)
}

func (r *relayAsEvcc) GetState() evcc.CarState {
	return r.state
}

func (r *relayAsEvcc) EnableCharging() error {
	if r.inverseLogic {
		_ = r.pin.SetValue(0)
	} else {
		_ = r.pin.SetValue(1)
	}

	return r.setState(evcc.StateC2, string(core.NoError))
}

func (r *relayAsEvcc) DisableCharging() {
	if r.inverseLogic {
		_ = r.pin.SetValue(1)
	} else {
		_ = r.pin.SetValue(0)
	}

	_ = r.setState(evcc.StateA1, string(core.NoError))
}

func (r *relayAsEvcc) setState(state evcc.CarState, error string) error {
	if !evcc.IsStateValid(state) {
		return nil
	}

	r.state = state
	r.statusChannel <- evcc.NewStateNotification(state, error)
	return nil
}

func (r *relayAsEvcc) SetMaxChargingCurrent(value float64) error {
	return nil
}

func (r *relayAsEvcc) GetMaxChargingCurrent() float64 {
	return 0.0
}

func (r *relayAsEvcc) Cleanup() error {
	return r.pin.Close()
}

func (r *relayAsEvcc) GetType() string {
	return Relay
}

func (r *relayAsEvcc) GetStatusChangeChannel() <-chan evcc.StateNotification {
	return r.statusChannel
}

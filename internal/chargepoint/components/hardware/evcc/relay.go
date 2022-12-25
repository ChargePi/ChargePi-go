package evcc

import (
	"context"
	"errors"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/warthog618/gpiod"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/evcc"
)

var (
	ErrInvalidPinNumber = errors.New("pin number must be greater than 0")
	ErrInitFailed       = errors.New("init failed")
)

type (
	RelayAsEvcc struct {
		relayPin      int
		inverseLogic  bool
		state         evcc.CarState
		pin           *gpiod.Line
		statusChannel chan chargePoint.StateNotification
	}
)

// NewRelay creates a new RelayImpl struct that will communicate with the GPIO pin specified.
func NewRelay(relayPin int, inverseLogic bool) (*RelayAsEvcc, error) {
	if relayPin <= 0 {
		return nil, ErrInvalidPinNumber
	}

	log.Debugf("Creating new relay at pin %d", relayPin)
	relay := RelayAsEvcc{
		relayPin:      relayPin,
		inverseLogic:  inverseLogic,
		statusChannel: make(chan chargePoint.StateNotification, 10),
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

func (r *RelayAsEvcc) GetState() evcc.CarState {
	return r.state
}

func (r *RelayAsEvcc) EnableCharging() error {
	if r.inverseLogic {
		_ = r.pin.SetValue(0)
	} else {
		_ = r.pin.SetValue(1)
	}

	_ = r.setState(evcc.StateB2, string(core.NoError))
	return r.setState(evcc.StateC2, string(core.NoError))
}

func (r *RelayAsEvcc) DisableCharging() {
	if r.inverseLogic {
		_ = r.pin.SetValue(1)
	} else {
		_ = r.pin.SetValue(0)
	}

	_ = r.setState(evcc.StateB1, string(core.NoError))
	_ = r.setState(evcc.StateA1, string(core.NoError))
}

func (r *RelayAsEvcc) setState(state evcc.CarState, error string) error {
	if !evcc.IsStateValid(state) {
		return nil
	}

	r.state = state
	r.statusChannel <- chargePoint.NewStateNotification(state, error)
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

func (r *RelayAsEvcc) GetStatusChangeChannel() <-chan chargePoint.StateNotification {
	return r.statusChannel
}

func (r *RelayAsEvcc) Reset() {
}

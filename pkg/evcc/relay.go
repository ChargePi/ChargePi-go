package evcc

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
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
	logger        *zap.Logger
}

// NewRelay creates a new RelayImpl struct that will communicate with the GPIO pin specified.
func NewRelay(logger *zap.Logger, settings *settings.Relay) (*RelayAsEvcc, error) {
	if settings.RelayPin <= 0 {
		return nil, ErrInvalidPinNumber
	}

	relay := RelayAsEvcc{
		logger:        logger.Named("relay"),
		relayPin:      settings.RelayPin,
		inverseLogic:  settings.InverseLogic,
		statusChannel: make(chan StateNotification, 10),
	}

	return &relay, nil
}

func (r *RelayAsEvcc) Init(ctx context.Context) error {
	r.logger.Debug("Initializing relay")
	// Refer to gpiod docs
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		return err
	}

	r.pin, err = c.RequestLine(r.relayPin, gpiod.AsOutput(0))
	return err
}

func (r *RelayAsEvcc) Lock() {
	r.logger.Debug("Locking relay")
}

func (r *RelayAsEvcc) Unlock() {
	r.logger.Debug("Unlocking relay")
}

func (r *RelayAsEvcc) GetError() string {
	return string(core.NoError)
}

func (r *RelayAsEvcc) GetState() CarState {
	return r.state
}

func (r *RelayAsEvcc) EnableCharging() error {
	r.logger.Debug("Enabling charging")
	if r.inverseLogic {
		_ = r.pin.SetValue(0)
	} else {
		_ = r.pin.SetValue(1)
	}

	_ = r.setState(StateB2, string(core.NoError))
	return r.setState(StateC2, string(core.NoError))
}

func (r *RelayAsEvcc) DisableCharging() {
	r.logger.Debug("Disabling charging")
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
	r.logger.Debug("Setting max charging current", zap.Float64("value", value))
	return nil
}

func (r *RelayAsEvcc) GetMaxChargingCurrent() float64 {
	return 0.0
}

func (r *RelayAsEvcc) Cleanup() error {
	r.logger.Debug("Cleaning up relay")
	return r.pin.Close()
}

func (r *RelayAsEvcc) GetType() string {
	return Relay
}

func (r *RelayAsEvcc) GetStatusChangeChannel() <-chan StateNotification {
	return r.statusChannel
}
func (r *RelayAsEvcc) SetNotificationChannel(notifications chan StateNotification) {
	r.statusChannel = notifications
}

func (r *RelayAsEvcc) Reset() {
	r.logger.Debug("Resetting relay")
}

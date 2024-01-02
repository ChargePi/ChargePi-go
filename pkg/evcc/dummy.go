package evcc

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type Dummy struct {
	logger        *log.Logger
	settings      *settings.EvccDummy
	notifications chan StateNotification
	currentState  CarState
	maxCurrent    float64
}

func NewDummy(settings *settings.EvccDummy) (*Dummy, error) {
	return &Dummy{
		settings:      settings,
		logger:        log.StandardLogger(),
		notifications: make(chan StateNotification),
	}, nil
}

func (d *Dummy) Init(ctx context.Context) error {
	d.logger.Info("Initializing dummy EVCC")
	return nil
}

func (d *Dummy) EnableCharging() error {
	d.logger.Info("Enabling charging")
	return nil
}

func (d *Dummy) DisableCharging() {
	d.logger.Info("Disabling charging")
}

func (d *Dummy) SetMaxChargingCurrent(value float64) error {
	d.logger.Infof("Setting max charging current to %f", value)
	d.maxCurrent = value
	return nil
}

func (d *Dummy) GetMaxChargingCurrent() float64 {
	d.logger.Info("Getting max charging current")
	return d.maxCurrent
}

func (d *Dummy) Lock() {
	d.logger.Info("Locking dummy connector")
}

func (d *Dummy) Unlock() {
	d.logger.Info("Unlocking dummy connector")
}

func (d *Dummy) GetState() CarState {
	d.logger.Info("Getting car state")
	return d.currentState
}

func (d *Dummy) GetError() string {
	d.logger.Info("Getting dummy error")
	return ""
}

func (d *Dummy) Cleanup() error {
	d.logger.Info("Cleaning up dummy")
	return nil
}

func (d *Dummy) GetType() string {
	return TypeDummy
}

func (d *Dummy) GetStatusChangeChannel() <-chan StateNotification {
	return d.notifications
}

func (d *Dummy) SetNotificationChannel(notifications chan StateNotification) {
	d.notifications = notifications
}

func (d *Dummy) Reset() {
	d.logger.Info("Resetting dummy")
}

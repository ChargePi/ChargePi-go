package evcc

import (
	"context"

	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

const (
	PhoenixEMCPPPETH = "EM-CP-PP-ETH"
	Relay            = "Relay"
	Western          = "Western"
)

type EVCC interface {
	Init(ctx context.Context) error
	EnableCharging() error
	DisableCharging()
	SetMaxChargingCurrent(value float64) error
	GetMaxChargingCurrent() float64
	Lock()
	Unlock()
	GetState() CarState
	GetError() string
	Cleanup() error
	GetType() string
	GetStatusChangeChannel() <-chan StateNotification
	SetNotificationChannel(notifications chan StateNotification)
	Reset()
	// SelfCheck() error
}

// NewEVCCFromType creates a new EVCC instance based on the provided type.
func NewEVCCFromType(evccSettings settings.EVCC) (EVCC, error) {
	switch evccSettings.Type {
	case Relay:
		return NewRelay(evccSettings.RelayPin, evccSettings.InverseLogic)
	case Western:
		return NewWesternController(1, evccSettings.Serial)
	default:
		return nil, nil
	}
}
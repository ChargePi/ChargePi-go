package evcc

import (
	"context"
	"errors"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/pkg/models/settings"
)

const (
	PhoenixEMCPPPETH = "EM-CP-PP-ETH"
	Relay            = "Relay"
	Western          = "Western"
	TypeDummy        = "Dummy"
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
	logger := zap.L()
	switch evccSettings.Type {
	case Relay:
		return NewRelay(logger, evccSettings.Relay)
	case TypeDummy:
		return NewDummy(logger, evccSettings.Dummy)
	default:
		return nil, errors.New("unsupported EVCC type: " + evccSettings.Type)
	}
}

package evcc

import (
	"context"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"

	"github.com/xBlaz3kx/ChargePi-go/pkg/models/evcc"
)

const (
	PhoenixEMCPPPETH = "EM-CP-PP-ETH"
	Relay            = "Relay"
	Western          = "Western"
)

type (
	EVCC interface {
		Init(ctx context.Context) error
		EnableCharging() error
		DisableCharging()
		SetMaxChargingCurrent(value float64) error
		GetMaxChargingCurrent() float64
		Lock()
		Unlock()
		GetState() evcc.CarState
		GetError() string
		Cleanup() error
		GetType() string
		GetStatusChangeChannel() <-chan StateNotification
		Reset()
	}
)

// NewEVCCFromType creates a new EVCC instance based on the provided type.
func NewEVCCFromType(evccSettings settings.EVCC) (EVCC, error) {
	switch evccSettings.Type {
	case Relay:
		return NewRelay(evccSettings.RelayPin, evccSettings.InverseLogic)
	default:
		return nil, nil
	}
}

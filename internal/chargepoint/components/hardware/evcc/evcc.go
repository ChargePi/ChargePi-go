package evcc

import (
	"context"

	evccLib "github.com/GLCharge/evcc-library-go"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
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
		GetStatusChangeChannel() <-chan chargePoint.StateNotification
		Reset()
	}
)

// NewPowerMeter creates a new power meter based on the connector settings.
func NewEVCCFromType(evccSettings settings.EVCC) (EVCC, error) {
	switch evccSettings.Type {
	case Relay:
		return NewRelay(evccSettings.RelayPin, evccSettings.InverseLogic)
	case Western:
		controller, err := evccLib.NewWesternController(1, evccSettings.Serial.DeviceAddress)
		if err != nil {
			return nil, err
		}

		return NewWesternController(controller), nil
	default:
		return nil, nil
	}
}

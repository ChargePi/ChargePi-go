package manager

import (
	"context"

	"github.com/ChargePi/ocppManager-go/ocpp_v16"

	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
)

type SettingsRepository interface {
	GetSettings(ctx context.Context) (*chargePoint.Settings, error)
	UpdateSettings(ctx context.Context, settings chargePoint.Settings) error
}

type OcppConfigurationRepository interface {
	GetOcppConfiguration(ctx context.Context, version int) (*ocpp_v16.Config, error)
	GeLatestOcppConfiguration(ctx context.Context) (*ocpp_v16.Config, error)
	SetOcppConfiguration(ctx context.Context, config ocpp_v16.Config) error
}

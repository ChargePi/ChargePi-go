package manager

import (
	"github.com/ChargePi/ocppManager-go/ocpp_v16"

	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
)

type SettingsRepository interface {
	GetSettings() (*chargePoint.Settings, error)
	UpdateSettings(settings chargePoint.Settings) error
}

type OcppConfigurationRepository interface {
	GetOcppConfiguration(version int) (*ocpp_v16.Config, error)
	GeLatestOcppConfiguration() (*ocpp_v16.Config, error)
	SetOcppConfiguration(ocpp_v16.Config) error
}

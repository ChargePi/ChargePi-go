package settings

import (
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type (
	EVSE struct {
		EvseId      int                 `json:"evseId,omitempty" yaml:"evseId" mapstructure:"evseId" validate:"required"`
		MaxPower    float32             `json:"maxPower" yaml:"maxPower" mapstructure:"maxPower" validate:"gt=0"`
		EVCC        settings.EVCC       `json:"evcc" yaml:"evcc" mapstructure:"evcc"`
		PowerMeter  settings.PowerMeter `json:"powerMeter" yaml:"powerMeter" mapstructure:"powerMeter"`
		Status      string              `json:"status,omitempty" yaml:"status" mapstructure:"status"`
		Reservation *int                `json:"reservation,omitempty" yaml:"reservation" mapstructure:"reservation"`
		Connectors  []Connector         `json:"connectors" yaml:"connectors" mapstructure:"connectors"`
	}

	Connector struct {
		ConnectorId int    `json:"connectorId,omitempty" yaml:"connectorId" mapstructure:"connectorId"`
		Type        string `json:"type,omitempty" yaml:"type" mapstructure:"type"`
		Status      string `json:"status,omitempty" yaml:"status" mapstructure:"status"`
	}
)

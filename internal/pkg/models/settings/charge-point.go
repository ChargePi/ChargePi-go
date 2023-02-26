package settings

import (
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
)

type (
	Settings struct {
		ChargePoint ChargePoint `json:"chargePoint" yaml:"chargePoint" mapstructure:"chargePoint"`
		Api         Api         `json:"api" yaml:"api" mapstructure:"api"`
	}

	ChargePoint struct {
		Info               Info               `json:"info" yaml:"info" mapstructure:"info"`
		ConnectionSettings ConnectionSettings `json:"connectionSettings" yaml:"connectionSettings" mapstructure:"connectionSettings"`
		Logging            Logging            `json:"logging" yaml:"logging" mapstructure:"logging"`
		Hardware           Hardware           `json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	}

	ConnectionSettings struct {
		Id                string               `validate:"required" json:"id,omitempty" yaml:"id" mapstructure:"id"`
		ProtocolVersion   ocpp.ProtocolVersion `validate:"required" json:"ProtocolVersion,omitempty" yaml:"ProtocolVersion" mapstructure:"ProtocolVersion"`
		ServerUri         string               `json:"ServerUri,omitempty" yaml:"ServerUri" mapstructure:"ServerUri"`
		BasicAuthUsername string               `json:"basicAuthUser,omitempty" yaml:"basicAuthUser" mapstructure:"basicAuthUser"`
		BasicAuthPassword string               `json:"basicAuthPass,omitempty" yaml:"basicAuthPass" mapstructure:"basicAuthPass"`
		TLS               TLS                  `json:"tls" yaml:"tls" mapstructure:"tls"`
	}

	Info struct {
		// Maximum time allowed if free mode is enabled
		MaxChargingTime *int `json:"MaxChargingTime,omitempty" yaml:"MaxChargingTime" mapstructure:"MaxChargingTime"`
		FreeMode        bool `json:"freeMode,omitempty" yaml:"freeMode" mapstructure:"freeMode"`
		// AC or DC
		Type string `json:"type,omitempty" yaml:"type" mapstructure:"type"`
		// in kW
		MaxPower    float32     `json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`
		OCPPDetails OCPPDetails `json:"ocpp" yaml:"ocpp" mapstructure:"ocpp"`
	}

	FreeChargingMode struct {
		Enabled  bool   `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Strategy string `json:"strategy,omitempty" yaml:"strategy" mapstructure:"strategy"`
	}
)

package settings

import (
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
)

type (
	Settings struct {
		ChargePoint ChargePoint `json:"chargePoint" yaml:"chargePoint" mapstructure:"chargePoint"`
		Api         Api         `json:"api" yaml:"api" mapstructure:"api"`
		Ui          Ui          `json:"ui" yaml:"ui" mapstructure:"ui"`
	}

	ChargePoint struct {
		Info               Info               `json:"info" yaml:"info" mapstructure:"info"`
		ConnectionSettings ConnectionSettings `json:"connectionSettings" yaml:"connectionSettings" mapstructure:"connectionSettings"`
		Logging            Logging            `json:"logging" yaml:"logging" mapstructure:"logging"`
		Hardware           Hardware           `json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	}

	ConnectionSettings struct {
		Id                string               `json:"id,omitempty" yaml:"id" mapstructure:"id" validate:"required"`
		ProtocolVersion   ocpp.ProtocolVersion `json:"protocolVersion,omitempty" yaml:"protocolVersion" mapstructure:"protocolVersion" validate:"required"`
		ServerUri         string               `json:"uri,omitempty" yaml:"uri" mapstructure:"uri" validate:"required"`
		BasicAuthUsername string               `json:"basicAuthUser,omitempty" yaml:"basicAuthUser" mapstructure:"basicAuthUser"`
		BasicAuthPassword string               `json:"basicAuthPass,omitempty" yaml:"basicAuthPass" mapstructure:"basicAuthPass"`
		TLS               TLS                  `json:"tls" yaml:"tls" mapstructure:"tls"`
	}

	Info struct {
		// Maximum time allowed if free mode is enabled
		MaxChargingTime *int `json:"MaxChargingTime,omitempty" yaml:"MaxChargingTime" mapstructure:"MaxChargingTime"`
		FreeMode        bool `json:"freeMode,omitempty" yaml:"freeMode" mapstructure:"freeMode"`
		// AC or DC
		Type string `json:"type,omitempty" yaml:"type" mapstructure:"type" validate:"oneof=AC DC"`
		// in kW
		MaxPower    float32     `json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`
		OCPPDetails OCPPDetails `json:"ocpp" yaml:"ocpp" mapstructure:"ocpp"`
	}

	FreeChargingMode struct {
		Enabled  bool   `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Strategy string `json:"strategy,omitempty" yaml:"strategy" mapstructure:"strategy"`
	}
)

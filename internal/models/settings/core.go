package settings

import "github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"

type (
	Settings struct {
		ChargePoint ChargePoint `fig:"chargePoint" json:"chargePoint" yaml:"chargePoint" mapstructure:"chargePoint"`
		Api         Api         `fig:"api" json:"api" yaml:"api" mapstructure:"api"`
	}

	ChargePoint struct {
		ConnectionSettings ConnectionSettings `fig:"connectionSettings" json:"connectionSettings" yaml:"connectionSettings" mapstructure:"connectionSettings"`
		Info               Info               `fig:"info" json:"info" yaml:"info" mapstructure:"info"`
		Logging            Logging            `fig:"logging" json:"logging" yaml:"logging" mapstructure:"logging"`
		Hardware           Hardware           `fig:"hardware" json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	}

	ConnectionSettings struct {
		Id                string               `fig:"Id" validate:"required" json:"id,omitempty" yaml:"id" mapstructure:"id"`
		ProtocolVersion   ocpp.ProtocolVersion `fig:"ProtocolVersion" default:"1.6" json:"ProtocolVersion,omitempty" yaml:"ProtocolVersion" mapstructure:"ProtocolVersion"`
		ServerUri         string               `fig:"ServerUri" validate:"required" json:"ServerUri,omitempty" yaml:"ServerUri" mapstructure:"ServerUri"`
		BasicAuthUsername string               `fig:"basicAuthUser" json:"basicAuthUser,omitempty" yaml:"basicAuthUser" mapstructure:"basicAuthUser"`
		BasicAuthPassword string               `fig:"basicAuthPass" json:"basicAuthPass,omitempty" yaml:"basicAuthPass" mapstructure:"basicAuthPass"`
		TLS               TLS                  `fig:"tls" json:"tls" yaml:"tls" mapstructure:"tls"`
	}

	Info struct {
		FreeChargingMode FreeChargingMode `fig:"freeChargingMode" default:"180" json:"freeChargingMode,omitempty" yaml:"freeChargingMode" mapstructure:"freeChargingMode"`
		OCPPInfo         OCPPInfo         `fig:"ocpp" json:"ocpp" yaml:"ocpp" mapstructure:"ocpp"`
	}

	FreeChargingMode struct {
		Enabled         bool `fig:"enabled" default:"180" json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		MaxChargingTime *int `fig:"MaxChargingTime" json:"MaxChargingTime,omitempty" yaml:"MaxChargingTime" mapstructure:"MaxChargingTime"`
	}
)

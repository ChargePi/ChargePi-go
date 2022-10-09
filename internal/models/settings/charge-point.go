package settings

type (
	Settings struct {
		ChargePoint ChargePoint `fig:"chargePoint" json:"chargePoint" yaml:"chargePoint" mapstructure:"chargePoint"`
		Api         Api         `fig:"api" json:"api" yaml:"api" mapstructure:"api"`
	}

	ChargePoint struct {
		Info               Info               `fig:"info" json:"info" yaml:"info" mapstructure:"info"`
		ConnectionSettings ConnectionSettings `fig:"connectionSettings" json:"connectionSettings" yaml:"connectionSettings" mapstructure:"connectionSettings"`
		Logging            Logging            `fig:"logging" json:"logging" yaml:"logging" mapstructure:"logging"`
		Hardware           Hardware           `fig:"hardware" json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	}

	ConnectionSettings struct {
		Id                string `fig:"Id" validate:"required" json:"id,omitempty" yaml:"id" mapstructure:"id"`
		ProtocolVersion   string `fig:"ProtocolVersion" default:"1.6" json:"ProtocolVersion,omitempty" yaml:"ProtocolVersion" mapstructure:"ProtocolVersion"`
		ServerUri         string `fig:"ServerUri" validate:"required" json:"ServerUri,omitempty" yaml:"ServerUri" mapstructure:"ServerUri"`
		BasicAuthUsername string `fig:"basicAuthUser" json:"basicAuthUser,omitempty" yaml:"basicAuthUser" mapstructure:"basicAuthUser"`
		BasicAuthPassword string `fig:"basicAuthPass" json:"basicAuthPass,omitempty" yaml:"basicAuthPass" mapstructure:"basicAuthPass"`
		TLS               TLS    `fig:"tls" json:"tls" yaml:"tls" mapstructure:"tls"`
	}

	Info struct {
		// Maximum time allowed if free mode is enabled
		MaxChargingTime int  `fig:"MaxChargingTime" default:"180" json:"MaxChargingTime,omitempty" yaml:"MaxChargingTime" mapstructure:"MaxChargingTime"`
		FreeMode        bool `fig:"freeMode" json:"freeMode,omitempty" yaml:"freeMode" mapstructure:"freeMode"`
		// AC or DC
		Type string `fig:"type" default:"180" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		// in kW
		MaxPower float32  `fig:"maxPower" default:"180" json:"maxPower,omitempty" yaml:"maxPower" mapstructure:"maxPower"`
		OCPPInfo OCPPInfo `fig:"ocpp" json:"ocpp" yaml:"ocpp" mapstructure:"ocpp"`
	}
)

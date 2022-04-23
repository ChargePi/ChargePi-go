package settings

type (
	Settings struct {
		ChargePoint ChargePoint `fig:"chargePoint" json:"chargePoint" yaml:"chargePoint" mapstructure:"chargePoint"`
		Api         Api         `fig:"api" json:"api" yaml:"api" mapstructure:"api"`
	}

	ChargePoint struct {
		Info     Info     `fig:"info" json:"info" yaml:"info" mapstructure:"info"`
		Logging  Logging  `fig:"logging" json:"logging" yaml:"logging" mapstructure:"logging"`
		TLS      TLS      `fig:"tls" json:"tls" yaml:"tls" mapstructure:"tls"`
		Hardware Hardware `fig:"hardware" json:"hardware" yaml:"hardware" mapstructure:"hardware"`
	}

	Info struct {
		Id              string   `fig:"Id" validate:"required" json:"id,omitempty" yaml:"id" mapstructure:"id"`
		ProtocolVersion string   `fig:"ProtocolVersion" default:"1.6" json:"ProtocolVersion,omitempty" yaml:"ProtocolVersion" mapstructure:"ProtocolVersion"`
		ServerUri       string   `fig:"ServerUri" validate:"required" json:"ServerUri,omitempty" yaml:"ServerUri" mapstructure:"ServerUri"`
		MaxChargingTime int      `fig:"MaxChargingTime" default:"180" json:"MaxChargingTime,omitempty" yaml:"MaxChargingTime" mapstructure:"MaxChargingTime"`
		OCPPInfo        OCPPInfo `fig:"ocpp" json:"ocpp" yaml:"ocpp" mapstructure:"ocpp"`
	}

	TLS struct {
		IsEnabled             bool   `fig:"isEnabled"  json:"isEnabled,omitempty" yaml:"isEnabled" mapstructure:"isEnabled"`
		CACertificatePath     string `fig:"CACertificatePath" json:"CACertificatePath,omitempty" yaml:"CACertificatePath" mapstructure:"CACertificatePath"`
		ClientCertificatePath string `fig:"ClientCertificatePath" json:"ClientCertificatePath,omitempty" yaml:"ClientCertificatePath" mapstructure:"ClientCertificatePath"`
		ClientKeyPath         string `fig:"ClientKeyPath" json:"ClientKeyPath,omitempty" yaml:"ClientKeyPath" mapstructure:"ClientKeyPath"`
	}

	Logging struct {
		Type   []string `fig:"type" validate:"required" json:"type,omitempty" yaml:"type" mapstructure:"type"`      // file, remote, console
		Format string   `fig:"format" default:"syslog" json:"format,omitempty" yaml:"format" mapstructure:"format"` // gelf, syslog, json, etc
		Host   string   `fig:"host" json:"host,omitempty" yaml:"host" mapstructure:"host"`
		Port   int      `fig:"port" default:"1514" json:"port,omitempty" yaml:"port" mapstructure:"port"`
	}

	Api struct {
		Enabled bool   `fig:"enabled" json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Address string `fig:"address" json:"address,omitempty" yaml:"address" mapstructure:"address"`
		Port    int    `fig:"port" json:"port,omitempty" yaml:"port" mapstructure:"port"`
	}
)

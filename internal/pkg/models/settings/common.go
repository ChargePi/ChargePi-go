package settings

type (
	TLS struct {
		IsEnabled             bool   `json:"isEnabled,omitempty" yaml:"isEnabled" mapstructure:"isEnabled"`
		CACertificatePath     string `json:"CACertificatePath,omitempty" yaml:"CACertificatePath" mapstructure:"CACertificatePath"`
		ClientCertificatePath string `json:"ClientCertificatePath,omitempty" yaml:"ClientCertificatePath" mapstructure:"ClientCertificatePath"`
		PrivateKeyPath        string `json:"PrivateKeyPath,omitempty" yaml:"PrivateKeyPath" mapstructure:"PrivateKeyPath"`
	}

	Logging struct {
		LogTypes []LogType `json:"logTypes,omitempty" yaml:"logTypes" mapstructure:"logTypes"`
	}

	LogType struct {
		Type    string  `json:"type,omitempty" yaml:"type" mapstructure:"type" validate:"required"` // file, remote, console
		Format  *string `json:"format,omitempty" yaml:"format" mapstructure:"format"`               // gelf, syslog, json, etc
		Address *string `json:"address,omitempty" yaml:"address" mapstructure:"address"`
	}

	Api struct {
		Enabled bool   `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Address string `validate:"hostname_port" json:"address,omitempty" yaml:"address" mapstructure:"address"`
		TLS     TLS    `json:"tls,omitempty" yaml:"tls" mapstructure:"tls"`
	}
)

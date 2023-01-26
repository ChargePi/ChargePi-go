package settings

type (
	TLS struct {
		IsEnabled             bool   `json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		CACertificatePath     string `json:"CACertificatePath,omitempty" yaml:"CACertificatePath" mapstructure:"CACertificatePath"`
		ClientCertificatePath string `json:"certificatePath,omitempty" yaml:"certificatePath" mapstructure:"certificatePath"`
		PrivateKeyPath        string `json:"keyPath,omitempty" yaml:"keyPath" mapstructure:"keyPath"`
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
		Address string `json:"address,omitempty" yaml:"address" mapstructure:"address" validate:"hostname_port"`
		TLS     TLS    `json:"tls,omitempty" yaml:"tls" mapstructure:"tls"`
	}
)

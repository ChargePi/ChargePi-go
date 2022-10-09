package settings

type (
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

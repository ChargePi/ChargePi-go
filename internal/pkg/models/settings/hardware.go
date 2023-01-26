package settings

type (
	Hardware struct {
		Display      Display      `json:"display" yaml:"display" mapstructure:"display"`
		TagReader    TagReader    `json:"tagReader" yaml:"tagReader" mapstructure:"tagReader"`
		LedIndicator LedIndicator `json:"ledIndicator" yaml:"ledIndicator" mapstructure:"ledIndicator"`
	}

	TagReader struct {
		IsEnabled   bool   `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`
		ReaderModel string `json:"readerModel,omitempty" yaml:"readerModel,omitempty" mapstructure:"readerModel,omitempty"`
		// Based on the type, get the connection details
		Device   *string `json:"deviceAddress,omitempty" yaml:"deviceAddress,omitempty" mapstructure:"deviceAddress,omitempty"`
		ResetPin *int    `json:"resetPin,omitempty" yaml:"resetPin,omitempty" mapstructure:"resetPin,omitempty"`
	}

	Display struct {
		IsEnabled bool   `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`
		Driver    string `json:"driver,omitempty" yaml:"driver,omitempty" mapstructure:"driver,omitempty"`
		Language  string `json:"language,omitempty" yaml:"language,omitempty" mapstructure:"language,omitempty"`
		// Based on the type, get the connection details
		I2C *I2C `fig:"i2c" json:"i2c,omitempty" yaml:"i2c,omitempty" mapstructure:"i2c,omitempty"`
	}
)

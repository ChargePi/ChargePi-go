package settings

type (
	/* ------------- Hardware structs ------------*/

	Hardware struct {
		Display      Display      `fig:"display" json:"display" yaml:"display" mapstructure:"display"`
		TagReader    TagReader    `fig:"tagReader" json:"tagReader" yaml:"tagReader" mapstructure:"tagReader"`
		LedIndicator LedIndicator `fig:"ledIndicator" json:"ledIndicator" yaml:"ledIndicator" mapstructure:"ledIndicator"`
	}

	TagReader struct {
		IsEnabled   bool   `fig:"Enabled" json:"Enabled,omitempty" yaml:"Enabled" mapstructure:"Enabled"`
		ReaderModel string `fig:"ReaderModel" json:"ReaderModel,omitempty" yaml:"ReaderModel" mapstructure:"ReaderModel"`
		// Based on the type, get the connection details
		Device   string `fig:"DeviceAddress" json:"deviceAddress,omitempty" yaml:"deviceAddress" mapstructure:"deviceAddress"`
		ResetPin int    `fig:"ResetPin" json:"ResetPin,omitempty" yaml:"ResetPin" mapstructure:"ResetPin"`
	}

	Display struct {
		IsEnabled bool   `fig:"Enabled" json:"Enabled,omitempty" yaml:"Enabled" mapstructure:"Enabled"`
		Driver    string `fig:"Driver" json:"driver,omitempty" yaml:"driver" mapstructure:"driver"`
		Language  string `fig:"Language" json:"language,omitempty" yaml:"language" mapstructure:"language"`
		// Based on the type, get the connection details
		I2C *I2C `fig:"i2c" json:"i2c,omitempty" yaml:"i2c" mapstructure:"i2c"`
	}
)

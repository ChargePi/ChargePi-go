package settings

type (
	/* ------------- Hardware structs ------------*/

	Hardware struct {
		Display      Display      `fig:"display" json:"display" yaml:"display" mapstructure:"display"`
		TagReader    TagReader    `fig:"tagReader" json:"tagReader" yaml:"tagReader" mapstructure:"tagReader"`
		LedIndicator LedIndicator `fig:"ledIndicator" json:"ledIndicator" yaml:"ledIndicator" mapstructure:"ledIndicator"`
	}

	EVCC struct {
		Type string
		// Based on the type, get the connection details
		RelayPin     int     `fig:"RelayPin" json:"RelayPin,omitempty" yaml:"RelayPin" mapstructure:"RelayPin"`
		InverseLogic bool    `fig:"InverseLogic" json:"InverseLogic,omitempty" yaml:"InverseLogic" mapstructure:"InverseLogic"`
		Serial       *Serial `fig:"serial" json:"serial,omitempty" yaml:"serial" mapstructure:"serial"`
		Modbus       *ModBus `fig:"modbus" json:"modbus,omitempty" yaml:"modbus" mapstructure:"modbus"`
	}

	LedIndicator struct {
		Enabled          bool   `fig:"Enabled" json:"Enabled,omitempty" yaml:"Enabled" mapstructure:"Enabled"`
		IndicateCardRead bool   `fig:"IndicateCardRead" json:"IndicateCardRead,omitempty" yaml:"IndicateCardRead" mapstructure:"IndicateCardRead"`
		Type             string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		// Based on the type, get the connection details
		DataPin int  `fig:"DataPin" json:"DataPin,omitempty" yaml:"DataPin" mapstructure:"DataPin"`
		Invert  bool `fig:"Invert" json:"invert,omitempty" yaml:"invert" mapstructure:"invert"`
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

	PowerMeter struct {
		Enabled bool   `fig:"Enabled" json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Type    string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		// Based on the type, get the connection details
		// For smarter energy meters, using Modbus RTU or similar based on the device type
		ModBus *ModBus `fig:"modbus" json:"modbus,omitempty" yaml:"modbus" mapstructure:"modbus"`
		// CS5460 specific details
		SPI                  *SPI     `fig:"spi" json:"spi,omitempty" yaml:"spi" mapstructure:"spi"`
		ShuntOffset          *float64 `fig:"ShuntOffset" json:"ShuntOffset,omitempty" yaml:"ShuntOffset" mapstructure:"ShuntOffset"`
		VoltageDividerOffset *float64 `fig:"VoltageDividerOffset" json:"VoltageDividerOffset,omitempty" yaml:"VoltageDividerOffset" mapstructure:"VoltageDividerOffset"`
	}
)

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
		// Based on the type
		RelayPin     int  `fig:"RelayPin" json:"RelayPin,omitempty" yaml:"RelayPin" mapstructure:"RelayPin"`
		InverseLogic bool `fig:"InverseLogic" json:"InverseLogic,omitempty" yaml:"InverseLogic" mapstructure:"InverseLogic"`
	}

	LedIndicator struct {
		Enabled          bool   `fig:"Enabled" json:"Enabled,omitempty" yaml:"Enabled" mapstructure:"Enabled"`
		DataPin          int    `fig:"DataPin" json:"DataPin,omitempty" yaml:"DataPin" mapstructure:"DataPin"`
		IndicateCardRead bool   `fig:"IndicateCardRead" json:"IndicateCardRead,omitempty" yaml:"IndicateCardRead" mapstructure:"IndicateCardRead"`
		Type             string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		Invert           bool   `fig:"Invert" json:"invert,omitempty" yaml:"invert" mapstructure:"invert"`
	}

	TagReader struct {
		IsEnabled   bool   `fig:"Enabled" json:"Enabled,omitempty" yaml:"Enabled" mapstructure:"Enabled"`
		ReaderModel string `fig:"ReaderModel" json:"ReaderModel,omitempty" yaml:"ReaderModel" mapstructure:"ReaderModel"`
		Device      string `fig:"DeviceAddress" json:"deviceAddress,omitempty" yaml:"deviceAddress" mapstructure:"deviceAddress"`
		ResetPin    int    `fig:"ResetPin" json:"ResetPin,omitempty" yaml:"ResetPin" mapstructure:"ResetPin"`
	}

	Display struct {
		IsEnabled  bool   `fig:"Enabled" json:"Enabled,omitempty" yaml:"Enabled" mapstructure:"Enabled"`
		Driver     string `fig:"Driver" json:"driver,omitempty" yaml:"driver" mapstructure:"driver"`
		Language   string `fig:"Language" json:"language,omitempty" yaml:"language" mapstructure:"language"`
		I2CAddress string `fig:"I2CAddress" json:"I2CAddress,omitempty" yaml:"I2CAddress" mapstructure:"I2CAddress"`
		I2CBus     int    `fig:"I2CBus" json:"I2CBus,omitempty" yaml:"I2CBus" mapstructure:"I2CBus"`
	}

	PowerMeter struct {
		Enabled bool   `fig:"Enabled" json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Type    string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`

		// For smarter energy meters, using Modbus RTU or similar based on the device type
		DeviceAddress string `fig:"deviceAddress" json:"deviceAddress,omitempty" yaml:"deviceAddress" mapstructure:"deviceAddress"`
		Protocol      string `fig:"protocol" json:"protocol,omitempty" yaml:"protocol" mapstructure:"protocol"`

		// For CS5460
		PowerMeterPin        int     `fig:"PowerMeterPin" json:"PowerMeterPin,omitempty" yaml:"PowerMeterPin" mapstructure:"PowerMeterPin"`
		SpiBus               int     `fig:"SpiBus" default:"0" json:"SpiBus,omitempty" yaml:"SpiBus" mapstructure:"SpiBus"`
		ShuntOffset          float64 `fig:"ShuntOffset" json:"ShuntOffset,omitempty" yaml:"ShuntOffset" mapstructure:"ShuntOffset"`
		VoltageDividerOffset float64 `fig:"VoltageDividerOffset" json:"VoltageDividerOffset,omitempty" yaml:"VoltageDividerOffset" mapstructure:"VoltageDividerOffset"`
	}
)

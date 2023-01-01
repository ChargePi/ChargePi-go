package settings

type (
	PowerMeter struct {
		Enabled bool   `fig:"Enabled" json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
		Type    string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		// Based on the type, get the connection details
		// For smarter energy meters, using Modbus RTU or similar based on the device type
		ModBus *ModBus `fig:"modbus" json:"modbus,omitempty" yaml:"modbus" mapstructure:"modbus"`
		SPI    *SPI    `fig:"spi" json:"spi,omitempty" yaml:"spi" mapstructure:"spi"`
		// CS5460 specific details
		CS5460 *CS5460 `fig:"CS5460" json:"CS5460,omitempty" yaml:"CS5460" mapstructure:"CS5460"`
	}

	CS5460 struct {
		ShuntOffset          float64 `fig:"ShuntOffset" json:"ShuntOffset,omitempty" yaml:"ShuntOffset" mapstructure:"ShuntOffset"`
		VoltageDividerOffset float64 `fig:"VoltageDividerOffset" json:"VoltageDividerOffset,omitempty" yaml:"VoltageDividerOffset" mapstructure:"VoltageDividerOffset"`
	}
)

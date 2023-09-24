package settings

type (
	PowerMeter struct {
		Enabled bool   `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
		Type    string `json:"type,omitempty" yaml:"type,omitempty" mapstructure:"type,omitempty"`
		// Based on the type, get the connection details
		// For smarter energy meters, using Modbus RTU or similar based on the device type
		ModBus *ModBus `json:"modbus,omitempty" yaml:"modbus,omitempty" mapstructure:"modbus,omitempty"`
		SPI    *SPI    `json:"spi,omitempty" yaml:"spi,omitempty" mapstructure:"spi,omitempty"`
		// CS5460 specific details
		CS5460 *CS5460 `json:"cs5460,omitempty" yaml:"cs5460,omitempty" mapstructure:"cs5460,omitempty"`
	}

	CS5460 struct {
		ShuntOffset          float64 `json:"shuntOffset,omitempty" yaml:"shuntOffset,omitempty" mapstructure:"shuntOffset,omitempty"`
		VoltageDividerOffset float64 `json:"voltageDividerOffset,omitempty" yaml:"voltageDividerOffset,omitempty" mapstructure:"voltageDividerOffset,omitempty"`
	}
)

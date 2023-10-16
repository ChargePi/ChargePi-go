package settings

type (
	Display struct {
		// Enable or disable the display from the configuration
		IsEnabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`

		// Enable the display to be controller remotely through OCPP
		RemoteEnabled bool `json:"remote,omitempty" yaml:"remote,omitempty" mapstructure:"remote,omitempty"`

		// Display driver/type of the display - can be a direct implementation of the driver or an HTTP server
		Driver string `json:"driver,omitempty" yaml:"driver,omitempty" mapstructure:"driver,omitempty"`

		// Default display language
		Language string `json:"language,omitempty" yaml:"language,omitempty" mapstructure:"language,omitempty"`

		// Hitachi HD44780 display configuration details
		HD44780 *HD44780 `json:"hd44780,omitempty" yaml:"hd44780,omitempty" mapstructure:"hd44780,omitempty"`

		// Dummy display configuration details
		DisplayDummy *DisplayDummy `json:"dummy,omitempty" yaml:"dummy,omitempty" mapstructure:"dummy,omitempty"`
	}

	HD44780 struct {
		I2C I2C `fig:"i2c" json:"i2c,omitempty" yaml:"i2c,omitempty" mapstructure:"i2c,omitempty"`
	}

	DisplayDummy struct {
	}
)

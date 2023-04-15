package settings

import (
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type (
	Hardware struct {
		Display   Display   `json:"display" yaml:"display" mapstructure:"display"`
		TagReader TagReader `json:"reader" yaml:"reader" mapstructure:"reader"`
		Indicator Indicator `json:"indicator" yaml:"indicator" mapstructure:"indicator"`
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
		I2C *settings.I2C `fig:"i2c" json:"i2c,omitempty" yaml:"i2c,omitempty" mapstructure:"i2c,omitempty"`
	}
)

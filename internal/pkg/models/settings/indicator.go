package settings

type (
	Indicator struct {
		Enabled           bool                    `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
		IndicateCardRead  bool                    `json:"indicateCardRead,omitempty" yaml:"indicateCardRead,omitempty" mapstructure:"indicateCardRead,omitempty"`
		Type              string                  `json:"type,omitempty" yaml:"type" mapstructure:"type"`
		IndicatorMappings *IndicatorStatusMapping `json:"statuses,omitempty" yaml:"statuses,omitempty" mapstructure:"statuses,omitempty"`
		// Based on the type, get the connection details
		WS281x *WS281x `json:"ws281x,omitempty" yaml:"ws281x,omitempty" mapstructure:"ws281x,omitempty"`
	}

	WS281x struct {
		DataPin int  `json:"dataPin,omitempty" yaml:"dataPin,omitempty" mapstructure:"dataPin,omitempty"`
		Invert  bool `json:"invert,omitempty" yaml:"invert,omitempty" mapstructure:"invert,omitempty"`
	}

	IndicatorStatusMapping struct {
		Available string `json:"available,omitempty" yaml:"available,omitempty" mapstructure:"available,omitempty" validate:"hexcolor"`
		Reserved  string `json:"reserved,omitempty" yaml:"reserved,omitempty" mapstructure:"reserved,omitempty" validate:"hexcolor"`
		Preparing string `json:"preparing,omitempty" yaml:"preparing,omitempty" mapstructure:"preparing,omitempty" validate:"hexcolor"`
		Charging  string `json:"charging,omitempty" yaml:"charging,omitempty" mapstructure:"charging,omitempty" validate:"hexcolor"`
		Finishing string `json:"finishing,omitempty" yaml:"finishing,omitempty" mapstructure:"finishing,omitempty" validate:"hexcolor"`
		Fault     string `json:"fault,omitempty" yaml:"fault,omitempty" mapstructure:"fault,omitempty" validate:"hexcolor"`
		Error     string `json:"error,omitempty" yaml:"error,omitempty" mapstructure:"error,omitempty" validate:"hexcolor"`
	}
)

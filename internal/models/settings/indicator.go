package settings

type (
	LedIndicator struct {
		Enabled          bool   `fig:"Enabled" json:"Enabled,omitempty" yaml:"Enabled" mapstructure:"Enabled"`
		IndicateCardRead bool   `fig:"IndicateCardRead" json:"IndicateCardRead,omitempty" yaml:"IndicateCardRead" mapstructure:"IndicateCardRead"`
		Type             string `fig:"Type" json:"type,omitempty" yaml:"type" mapstructure:"type"`
		// Based on the type, get the connection details
		DataPin           int                    `fig:"DataPin" json:"DataPin,omitempty" yaml:"DataPin" mapstructure:"DataPin"`
		Invert            bool                   `fig:"Invert" json:"invert,omitempty" yaml:"invert" mapstructure:"invert"`
		IndicatorMappings IndicatorStatusMapping `fig:"statuses" json:"statuses,omitempty" yaml:"statuses" mapstructure:"statuses"`
	}

	IndicatorStatusMapping struct {
		Available string `json:"available,omitempty" yaml:"available" mapstructure:"available"`
		Reserved  string `json:"reserved,omitempty" yaml:"reserved" mapstructure:"reserved"`
		Preparing string `json:"preparing,omitempty" yaml:"preparing" mapstructure:"preparing"`
		Charging  string `json:"charging,omitempty" yaml:"charging" mapstructure:"charging"`
		Finishing string `json:"finishing,omitempty" yaml:"finishing" mapstructure:"finishing"`
		Fault     string `json:"fault,omitempty" yaml:"fault" mapstructure:"fault"`
		Error     string `json:"error,omitempty" yaml:"error" mapstructure:"error"`
	}
)

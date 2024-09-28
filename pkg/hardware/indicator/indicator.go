package indicator

// Indicator is an abstraction layer for connector status indication, usually an RGB LED strip.
type Indicator interface {
	ChangeColor(index int, color Color) error
	Blink(index int, times int, color Color) error
	SetBrightness(brightness int) error
	GetBrightness() int
	Cleanup()
	GetType() string
}

type StatusMapping struct {
	Available string `json:"available,omitempty" yaml:"available,omitempty" mapstructure:"available,omitempty" validate:"hexcolor"`
	Reserved  string `json:"reserved,omitempty" yaml:"reserved,omitempty" mapstructure:"reserved,omitempty" validate:"hexcolor"`
	Preparing string `json:"preparing,omitempty" yaml:"preparing,omitempty" mapstructure:"preparing,omitempty" validate:"hexcolor"`
	Charging  string `json:"charging,omitempty" yaml:"charging,omitempty" mapstructure:"charging,omitempty" validate:"hexcolor"`
	Finishing string `json:"finishing,omitempty" yaml:"finishing,omitempty" mapstructure:"finishing,omitempty" validate:"hexcolor"`
	Fault     string `json:"fault,omitempty" yaml:"fault,omitempty" mapstructure:"fault,omitempty" validate:"hexcolor"`
	Error     string `json:"error,omitempty" yaml:"error,omitempty" mapstructure:"error,omitempty" validate:"hexcolor"`
}

type Color string

// color constants
const (
	Off    = Color("Off")
	White  = Color("White")
	Red    = Color("Red")
	Green  = Color("Green")
	Blue   = Color("Blue")
	Yellow = Color("Yellow")
	Orange = Color("Orange")
)

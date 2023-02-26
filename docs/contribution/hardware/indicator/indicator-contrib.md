## Indicator hardware

All indicators should implement the `Indicator` interface. Initialize all the communication when creating a new struct.

```golang
package indicator

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

// Supported types
const (
	TypeWS281x = "WS281x"
)

type (
	Color string

	// Indicator is an abstraction layer for connector status indication, usually an RGB LED strip.
	Indicator interface {
		DisplayColor(index int, color Color) error
		Blink(index int, times int, color Color) error
		Cleanup()
		GetType() string
	}
)

```
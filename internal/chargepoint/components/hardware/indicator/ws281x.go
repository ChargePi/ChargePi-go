//go:build cgo
// +build cgo

package indicator

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"time"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	ocppManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

const (
	brightness = 128
	freq       = 800000
	sleepTime  = 500
)

// color constants
const (
	wS281xOff    = uint32(0x0)
	wS281xWhite  = uint32(0xFFFFFF)
	wS281xRed    = uint32(0xff0000)
	wS281xGreen  = uint32(0x00ff00)
	wS281xBlue   = uint32(0x000ff)
	wS281xYellow = uint32(0xeeff00)
	wS281xOrange = uint32(0xff7b00)
)

type WS281x struct {
	numberOfLEDs int
	dataPin      int
	ws2811       *ws2811.WS2811
}

// NewWS281xStrip create a new LED strip object with the specified number of LEDs and the data pin.
// When created, it will also be initialized.
func NewWS281xStrip(numberOfLEDs int, dataPin int) (*WS281x, error) {
	if numberOfLEDs <= 0 {
		return nil, ErrInvalidNumberOfLeds
	}

	if dataPin <= 0 {
		return nil, ErrInvalidPin
	}

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = numberOfLEDs
	opt.Channels[0].GpioPin = dataPin
	opt.Frequency = freq

	setLightIntensity(opt)

	// Create a new strip
	ledStrip, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return nil, err
	}

	// Initialize the strip
	err = ledStrip.Init()
	if err != nil {
		return nil, err
	}

	return &WS281x{
		dataPin:      dataPin,
		numberOfLEDs: numberOfLEDs,
		ws2811:       ledStrip,
	}, nil
}

// DisplayColor change the color of the LED at specified index to the specified color.
// The index must be greater than 0 and less than the length of the LED strip.
func (ws *WS281x) DisplayColor(index int, color Color) error {
	if index < 0 || index > len(ws.ws2811.Leds(0)) {
		return ErrInvalidIndex
	}

	ws.ws2811.Leds(0)[index] = getColorAsHex(color)
	return ws.ws2811.Render()
}

// Blink the LED at index a certain number of times with the specified color. If the number of times the LED is supposed to blink is even, it will stay turned off after the blinking,
// otherwise it will stay on after the blinking.
func (ws *WS281x) Blink(index int, times int, color Color) error {
	if index < 0 || index > len(ws.ws2811.Leds(0)) {
		return ErrInvalidIndex
	}

	for i := 0; i < times; i++ {
		if i%2 == 0 {
			ws.ws2811.Leds(0)[index] = wS281xOff
		} else {
			ws.ws2811.Leds(0)[index] = getColorAsHex(color)
		}

		err := ws.ws2811.Render()
		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond * sleepTime)
	}

	return nil
}

// Cleanup turn the LEDs off and terminate the data connection.
func (ws *WS281x) Cleanup() {
	var i = 0

	for ws.numberOfLEDs != i {
		_ = ws.DisplayColor(i, Off)
		i++
	}

	ws.close()
}

func (ws *WS281x) close() {
	ws.ws2811.Fini()
}

func (ws *WS281x) GetType() string {
	return TypeWS281x
}

func getColorAsHex(color Color) uint32 {
	switch color {
	case White:
		return wS281xWhite
	case Off:
		return wS281xOff
	case Blue:
		return wS281xBlue
	case Green:
		return wS281xGreen
	case Orange:
		return wS281xOrange
	case Red:
		return wS281xRed
	case Yellow:
		return wS281xYellow
	default:
		customColor, err := hex.DecodeString(string(color))
		if err != nil {
			return wS281xOff
		}

		return binary.LittleEndian.Uint32(customColor)
	}
}

func setLightIntensity(opts ws2811.Option) {
	lightIntensity, confErr := ocppManager.GetConfigurationValue(configuration.LightIntensity.String())
	if confErr == nil {
		intensity, err := strconv.ParseFloat(*lightIntensity, 32)
		if err != nil {
			return
		}

		// Light intensity is in percent
		opts.Channels[0].Brightness = int(intensity / 100 * 255)
	}
}

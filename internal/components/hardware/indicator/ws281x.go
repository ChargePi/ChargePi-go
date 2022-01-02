package indicator

import (
	"github.com/rpi-ws281x/rpi-ws281x-go"
	"time"
)

const (
	brightness = 128
	freq       = 800000
	sleepTime  = 500
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

	ledStrip := &WS281x{dataPin: dataPin, numberOfLEDs: numberOfLEDs, ws2811: nil}
	err := ledStrip.init()
	if err != nil {
		return nil, err
	}

	return ledStrip, nil
}

// init initialize the LED strip.
func (ws *WS281x) init() error {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ws.numberOfLEDs
	opt.Channels[0].GpioPin = ws.dataPin
	opt.Frequency = freq

	ledStrip, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		return err
	}

	ws.ws2811 = ledStrip
	return ws.ws2811.Init()
}

// DisplayColor change the color of the LED at specified index to the specified color.
// The index must be greater than 0 and less than the length of the LED strip.
func (ws *WS281x) DisplayColor(index int, colorHex uint32) error {
	if index < 0 || index > len(ws.ws2811.Leds(0)) {
		return ErrInvalidIndex
	}

	ws.ws2811.Leds(0)[index] = colorHex
	return ws.ws2811.Render()
}

// Blink the LED at index a certain number of times with the specified color. If the number of times the LED is supposed to blink is even, it will stay turned off after the blinking,
// otherwise it will stay on after the blinking.
func (ws *WS281x) Blink(index int, times int, colorHex uint32) error {
	if index < 0 || index > len(ws.ws2811.Leds(0)) {
		return ErrInvalidIndex
	}

	for i := 0; i < times; i++ {
		if i%2 == 0 {
			ws.ws2811.Leds(0)[index] = Off
		} else {
			ws.ws2811.Leds(0)[index] = colorHex
		}
		ws.ws2811.Render()
		time.Sleep(time.Millisecond * sleepTime)
	}

	return nil
}

// Cleanup turn the LEDs off and terminate the data connection.
func (ws WS281x) Cleanup() {
	var i = 0
	for ws.numberOfLEDs != i {
		ws.DisplayColor(i, Off)
	}
	ws.close()
}

func (ws *WS281x) close() {
	ws.ws2811.Fini()
}

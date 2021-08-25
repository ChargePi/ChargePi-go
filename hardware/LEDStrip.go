package hardware

import (
	"errors"
	"github.com/rpi-ws281x/rpi-ws281x-go"
	"time"
)

const (
	brightness = 128
	freq       = 800000
	sleepTime  = 500
	OFF        = 0x0
	WHITE      = 0xFFFFFF
	RED        = 0xff0000
	GREEN      = 0x00ff00
	BLUE       = 0x000ff
	YELLOW     = 0xeeff00
	ORANGE     = 0xff7b00
)

type (
	LEDStrip struct {
		numberOfLEDs int
		dataPin      int
		ws2811       *ws2811.WS2811
	}
)

// NewLEDStrip create a new LED strip object with the specified number of LEDs and the data pin.
// When created, it will also be initialized.
func NewLEDStrip(numberOfLEDs int, dataPin int) (*LEDStrip, error) {
	if numberOfLEDs < 0 {
		return nil, errors.New("number of leds must be greater than zero")
	}
	if dataPin < 0 {
		return nil, errors.New("number of leds must be greater than zero")
	}
	ledStrip := &LEDStrip{dataPin: dataPin, numberOfLEDs: numberOfLEDs, ws2811: nil}
	err := ledStrip.init()
	if err != nil {
		return nil, err
	}
	return ledStrip, nil
}

// init initialize the LED strip.
func (ws *LEDStrip) init() error {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ws.numberOfLEDs
	opt.Channels[0].GpioPin = ws.dataPin
	opt.Frequency = freq
	ledStrip, err := ws2811.MakeWS2811(&opt)
	ws.ws2811 = ledStrip
	err = ws.ws2811.Init()
	if err != nil {
		return err
	}
	return nil
}

// DisplayColor change the color of the LED at specified index to the specified color.
// The index must be greater than 0 and less than the length of the LED strip.
func (ws *LEDStrip) DisplayColor(index int, colorHex uint32) error {
	if index < 0 || index > len(ws.ws2811.Leds(0)) {
		return errors.New("invalid index")
	}
	ws.ws2811.Leds(0)[index] = colorHex
	return ws.ws2811.Render()
}

// Blink blink the LED at index a certain number of times with the specified color. If the number of times the LED is supposed to blink is even, it will stay turned off after the blinking,
// otherwise it will stay on after the blinking.
func (ws *LEDStrip) Blink(index int, times int, colorHex uint32) error {
	if index < 0 || index > len(ws.ws2811.Leds(0)) {
		return errors.New("invalid index")
	}
	for i := 0; i < times; i++ {
		if i%2 == 0 {
			ws.ws2811.Leds(0)[index] = OFF
		} else {
			ws.ws2811.Leds(0)[index] = colorHex
		}
		ws.ws2811.Render()
		time.Sleep(time.Millisecond * sleepTime)
	}
	return nil
}

// Cleanup turn the LEDs off and terminate the data connection.
func (ws LEDStrip) Cleanup() {
	var i = 0
	for ws.numberOfLEDs != i {
		_ = ws.DisplayColor(i, OFF)
	}
	ws.close()
}

func (ws *LEDStrip) close() {
	ws.ws2811.Fini()
}

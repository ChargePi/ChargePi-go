package display

import (
	"github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"strconv"
	"time"
)

type HD44780 struct {
	LCDChannel chan chargePoint.Message
	i2c        *i2c.I2C
	display    *hd44780.Lcd
}

// NewHD44780 Create a new HD44780 struct.
func NewHD44780(i2cAddress string, i2cBus int) (*HD44780, error) {
	var display = HD44780{}

	decodeString, err := strconv.ParseUint(i2cAddress, 16, 8)
	if err != nil {
		return nil, err
	}

	i2cDev, err := i2c.NewI2C(uint8(decodeString), i2cBus)
	if err != nil {
		return nil, err
	}

	display.i2c = i2cDev

	// Construct the display with I2C connection
	lcd2, err := hd44780.NewLcd(display.i2c, hd44780.LCD_16x2)
	if err != nil {
		return nil, err
	}

	display.display = lcd2
	_ = lcd2.BacklightOn()
	_ = lcd2.Clear()
	return &display, nil
}

// DisplayMessage displays the message on the Display. Pairs of messages will be displayed for the duration set in Message.
func (lcd *HD44780) DisplayMessage(message chargePoint.Message) {
	log.Debugf("Displaying the message to Display: %v", message.Messages)

	// Display lines in pairs. If there are odd number of lines, display the last line by itself.
	for i := 0; i < len(message.Messages); i = i + 2 {
		_ = lcd.display.Clear()
		_ = lcd.display.ShowMessage(message.Messages[i], hd44780.SHOW_LINE_1)

		// Prevents index-out-of-range error
		if i < len(message.Messages)-1 {
			_ = lcd.display.ShowMessage(message.Messages[i+1], hd44780.SHOW_LINE_2)
		}

		time.Sleep(message.MessageDuration)
	}
}

func (lcd *HD44780) Clear() {
	_ = lcd.display.Clear()
}

// Cleanup Close the Display I2C connection.
func (lcd *HD44780) Cleanup() {
	close(lcd.LCDChannel)
	lcd.Clear()
	_ = lcd.display.BacklightOff()
	_ = lcd.i2c.Close()
}

func (lcd *HD44780) GetType() string {
	return DriverHD44780
}

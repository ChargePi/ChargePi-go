package display

import (
	"context"
	"github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
	log "github.com/sirupsen/logrus"
	"time"
)

type HD44780 struct {
	LCDChannel chan LCDMessage
	i2c        *i2c.I2C
	display    *hd44780.Lcd
}

// NewHD44780 Create a new HD44780 struct.
func NewHD44780(lcdChannel chan LCDMessage, i2cAddress string, i2cBus int) (*HD44780, error) {
	var display = HD44780{LCDChannel: lcdChannel}

	//todo fix i2c address resolution
	i2cDev, err := i2c.NewI2C(0x27, i2cBus)
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
	lcd2.BacklightOn()
	lcd2.Clear()
	return &display, nil
}

// DisplayMessage displays the message on the LCD. Pairs of messages will be displayed for the duration set in LCDMessage.
func (lcd *HD44780) DisplayMessage(message LCDMessage) {
	log.Debugf("Displaying the message to LCD: %v", message.Messages)

	// Display lines in pairs. If there are odd number of lines, display the last line by itself.
	for i := 0; i < len(message.Messages); i = i + 2 {
		lcd.display.Clear()
		lcd.display.ShowMessage(message.Messages[i], hd44780.SHOW_LINE_1)

		// Prevents index-out-of-range error
		if i < len(message.Messages)-1 {
			lcd.display.ShowMessage(message.Messages[i+1], hd44780.SHOW_LINE_2)
		}
		time.Sleep(message.MessageDuration)
	}
}

//ListenForMessages Listen for incoming message requests and display the message received.
func (lcd *HD44780) ListenForMessages(ctx context.Context) {
Listener:
	for {
		select {
		case <-ctx.Done():
			lcd.Cleanup()
			break Listener
		case message := <-lcd.LCDChannel:
			lcd.DisplayMessage(message)
			break
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (lcd HD44780) GetLcdChannel() chan<- LCDMessage {
	return lcd.LCDChannel
}

func (lcd HD44780) Clear() {
	lcd.display.Clear()
}

// Cleanup Close the LCD I2C connection.
func (lcd *HD44780) Cleanup() {
	close(lcd.LCDChannel)
	lcd.Clear()
	lcd.display.BacklightOff()
	lcd.i2c.Close()
}

package hardware

import (
	"fmt"
	"github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
	"log"
	"time"
)

// LCDMessage Object representing the message that will be displayed on the LCD.
// Each array element in Messages represents a line being displayed on the 16x2 screen.
type LCDMessage struct {
	Messages        []string
	messageDuration int
}

type LCD struct {
	LCDChannel chan LCDMessage
	i2c        *i2c.I2C
	Display    *hd44780.Lcd
}

// NewLCD Create a new LCD object.
func NewLCD(lcdChannel chan LCDMessage) *LCD {
	var display = LCD{LCDChannel: lcdChannel}
	i2cDev, err := i2c.NewI2C(0x27, 1)
	if err != nil {
		log.Fatal(err)
	}
	display.i2c = i2cDev
	// Construct lcd-device connected via I2C connection
	lcd2, err := hd44780.NewLcd(display.i2c, hd44780.LCD_16x2)
	if err != nil {
		log.Fatal(err)
	}
	display.Display = lcd2
	err = lcd2.BacklightOn()
	return &display
}

// DisplayMessage Display the message on the LCD. Pairs of messages will be displayed for the duration set in LCDMessage.
func (lcd *LCD) DisplayMessage(message LCDMessage) {
	for i := 0; i < len(message.Messages)-1; i++ {
		if i%2 == 1 {
			lcd.Display.ShowMessage(message.Messages[i], hd44780.SHOW_LINE_1)
			lcd.Display.ShowMessage(message.Messages[i+1], hd44780.SHOW_LINE_2)
			time.Sleep(time.Duration(message.messageDuration) * time.Second)
		}
	}
}

// Cleanup Close the LCD I2C connection.
func (lcd *LCD) Cleanup() {
	lcd.Display.BacklightOff()
	lcd.i2c.Close()
}

//DisplayMessages Listen for incoming message requests and display the message received.
func (lcd *LCD) DisplayMessages() {
	for {
		select {
		case message := <-lcd.LCDChannel:
			{
				fmt.Printf("LCD displays message:%s", message)
				lcd.DisplayMessage(message)
				break
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}

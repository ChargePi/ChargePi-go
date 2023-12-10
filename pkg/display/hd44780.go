package display

import (
	"strconv"
	"time"

	"github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
	"github.com/go-co-op/gocron"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type HD44780 struct {
	i2c       *i2c.I2C
	display   *hd44780.Lcd
	scheduler *gocron.Scheduler
}

// NewHD44780 Create a new HD44780 struct.
func NewHD44780(settings settings.HD44780) (*HD44780, error) {
	// Decode the I2C address from hex to uint8
	decodeString, err := strconv.ParseUint(settings.I2C.Address, 16, 8)
	if err != nil {
		return nil, err
	}

	// Establish I2C connection
	i2cDev, err := i2c.NewI2C(uint8(decodeString), settings.I2C.Bus)
	if err != nil {
		return nil, err
	}

	// Construct the display with I2C connection
	lcd2, err := hd44780.NewLcd(i2cDev, hd44780.LCD_16x2)
	if err != nil {
		return nil, err
	}

	_ = lcd2.BacklightOn()
	_ = lcd2.Clear()

	return &HD44780{
		i2c:       i2cDev,
		display:   lcd2,
		scheduler: scheduler.NewScheduler(),
	}, nil
}

// DisplayMessage displays the message on the Display. Pairs of messages will be displayed for the duration set in Message.
func (lcd *HD44780) DisplayMessage(message display.MessageInfo) {

	// Schedule the display of the message a bit later than the start time to prevent the message from being displayed
	if message.StartDateTime != nil {
		_, err := lcd.scheduler.At(*message.StartDateTime).Tag("displayMessage").Do(lcd.DisplayMessage, message)
		if err != nil {
			log.WithError(err).Errorf("Error scheduling ClearMessage")
		}
		return
	}

	log.Debugf("Displaying the message to Display: %v", message)

	// Display lines in pairs. If there are odd number of lines, display the last line by itself.
	lines := splitString(message.Message.Content, 16)
	for i := 0; i < len(lines); i = i + 2 {
		_ = lcd.display.Clear()
		_ = lcd.display.ShowMessage(lines[i], hd44780.SHOW_LINE_1)

		// Prevents index-out-of-range error
		if i < len(lines)-1 {
			_ = lcd.display.ShowMessage(lines[i+1], hd44780.SHOW_LINE_2)
		}

		time.Sleep(time.Second * 5)
	}

	if message.EndDateTime != nil {
		_, err := lcd.scheduler.At(*message.EndDateTime).Tag("clearMessage").Do(lcd.Clear)
		if err != nil {
			log.WithError(err).Errorf("Error scheduling ClearMessage")
		}
	}
}

func splitString(str string, size int) []string {
	var result []string

	for start := 0; start < len(str); start += size {
		end := start + size
		if end > len(str) {
			end = len(str)
		}
		result = append(result, str[start:end])
	}

	return result
}

func (lcd *HD44780) Clear() {
	_ = lcd.display.Clear()
}

// Cleanup Close the Display I2C connection.
func (lcd *HD44780) Cleanup() {
	lcd.Clear()
	_ = lcd.display.BacklightOff()
	_ = lcd.i2c.Close()
}

func (lcd *HD44780) GetType() string {
	return DriverHD44780
}

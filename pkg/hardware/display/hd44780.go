//go:build cgo && raspberrypi4

package display

import (
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/d2r2/go-hd44780"
	"github.com/d2r2/go-i2c"
	"github.com/go-co-op/gocron"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"

	"github.com/ChargePi/ChargePi-go/internal/pkg/scheduler"
	"github.com/ChargePi/ChargePi-go/pkg/hardware"
)

const DriverHD44780 = "hd44780"

type HD44780Settings struct {
	I2C hardware.I2C `fig:"i2c" json:"i2c,omitempty" yaml:"i2c,omitempty" mapstructure:"i2c,omitempty"`
}

type HD44780 struct {
	i2c       *i2c.I2C
	display   *hd44780.Lcd
	scheduler *gocron.Scheduler
	logger    *zap.Logger
}

// NewHD44780 Create a new HD44780 struct.
func NewHD44780(logger *zap.Logger, settings HD44780Settings) (*HD44780, error) {
	logger := log.StandardLogger()
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

	err = lcd2.BacklightOn()
	if err != nil {
		// Log the error but continue
		logger.WithError(err).Error("Error turning on backlight")
	}

	err = lcd2.Clear()
	if err != nil {
		// Log the error but continue
		logger.WithError(err).Error("Error clearing display")
		return nil, err
	}

	return &HD44780{
		logger:    logger.Named("hd44780_display"),
		i2c:       i2cDev,
		display:   lcd2,
		scheduler: scheduler.NewScheduler(),
		logger:    logger.WithField("component", "display-hd44780"),
	}, nil
}

// DisplayMessage displays the message on the Display. Pairs of messages will be displayed for the duration set in Message.
func (lcd *HD44780) DisplayMessage(message display.MessageInfo) error {
	lcd.logger.Sugar().Debugf("Displaying the message to Settings: %v", message)

	// Display lines in pairs. If there are odd number of lines, display the last line by itself.
	lines := splitStringToLines(message.Message.Content, 16)
	for i := 0; i < len(lines); i = i + 2 {
		_ = lcd.display.Clear()
		_ = lcd.display.ShowMessage(lines[i], hd44780.SHOW_LINE_1)

		// Prevents index-out-of-range error
		if i < len(lines)-1 {
			_ = lcd.display.ShowMessage(lines[i+1], hd44780.SHOW_LINE_2)
		}

		time.Sleep(time.Second * 5)
	}
}

func (lcd *HD44780) Clear() error {
	lcd.logger.Info("Clearing display")

	err := lcd.display.Clear()
	if err != nil {
		lcd.logger.WithError(err).Error("Error clearing display")
		return err
	}

	return nil
}

// Cleanup Close the Display I2C connection.
func (lcd *HD44780) Cleanup(ctx context.Context) error {
	lcd.logger.Info("Cleaning up display")

	lcd.Clear()
	err := lcd.display.BacklightOff()
	if err != nil {
		lcd.logger.WithError(err).Error("Error turning off backlight")
	}

	err = lcd.i2c.Close()
	if err != nil {
		lcd.logger.WithError(err).Error("Error closing I2C connection")
		return err
	}
}

func (lcd *HD44780) GetType() string {
	return DriverHD44780
}

func splitStringToLines(str string, size int) []string {
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

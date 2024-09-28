package v16

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"strings"
	"time"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/reader"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	displaMessages "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
)

func (cp *ChargePoint) SetReader(reader reader.Reader) error {
	if util.IsNilInterfaceOrPointer(reader) {
		return nil
	}

	cp.logger.Debug("Setting reader")
	cp.tagReader = reader
	return nil
}

func (cp *ChargePoint) SetDisplay(display display.Display) error {
	if util.IsNilInterfaceOrPointer(display) {
		return nil
	}

	cp.logger.Debug("Setting display")
	cp.display = display

	err := cp.display.Clear()
	if err != nil {
		return err
	}

	return nil
}

func (cp *ChargePoint) SetIndicator(indicator indicator.Indicator) error {
	if util.IsNilInterfaceOrPointer(indicator) {
		return nil
	}

	cp.logger.Debug("Setting indicator")
	cp.indicator = indicator
	return nil
}

// DisplayMessage displays a message without blocking the caller.
func (cp *ChargePoint) DisplayMessage(message displaMessages.MessageInfo) error {
	if util.IsNilInterfaceOrPointer(cp.display) {
		cp.logger.Warn("Cannot send message to display, it is disabled or not configured")
		return errors.New("display not configured")
	}

	logger := cp.logger.With(zap.Any("message",message))

	logger.Debug("Displaying message")
	go func() {
		err := cp.display.DisplayMessage(message)
		if err != nil {
			cp.logger.With(zap.Error(err)).Error("Error displaying message")
		}
	}()

	cp.logger.Sugar().Debugf("Sending message to display: %v", message)
	go cp.display.DisplayMessage(message)
	return nil
}

// indicateStatusChange Indicates a status change on the indicator.
func (cp *ChargePoint) indicateStatusChange(connectorIndex int, status core.ChargePointStatus) {
	logger := cp.logger.With(zap.Int("connector", connectorIndex+1))
	if util.IsNilInterfaceOrPointer(cp.indicator) {
		logger.Warn("Cannot indicate status change, indicator disabled or not configured")
		return
	}

	// Get the color for the status
	color, err := colorMapping(cp.indicatorMapping, status)
	if err != nil {
		logger.With(zap.Error(err)).Error("Error indicating status")
	}

	logger.Sugar().Debugf("Indicating connector status: %x", color)
	err = cp.indicator.ChangeColor(connectorIndex, *color)
	if err != nil {
		logger.With(zap.Error(err)).Error("Error indicating status")
	}
}

// indicateCardRead Blinks the LED to indicate that the card was read.
func (cp *ChargePoint) indicateCardRead(index int, color indicator.Color) {
	logger := cp.logger.With(zap.Int("connector", index+1))
	if util.IsNilInterfaceOrPointer(cp.indicator) {
		logger.Warn("Cannot indicate card read, disabled or not configured")
		return
	}

	logger.Debug("Indicating a tag was read")
	err := cp.indicator.Blink(index, 3, color)
	if err != nil {
		logger.With(zap.Error(err)).Error("Could not indicate a tag was read")
	}
}

// ListenForTag Listen for an RFID/NFC tag in a separate goroutine.
// If a tag is detected, it will return the tag ID, otherwise it will return an error.
func (cp *ChargePoint) ListenForTag(ctx context.Context) (*string, error) {
	if util.IsNilInterfaceOrPointer(cp.tagReader) {
		return nil, errors.New("reader not configured")
	}

	tagChannel := cp.tagReader.GetTagChannel()

	if util.IsNilInterfaceOrPointer(tagChannel) {
		return nil, errors.New("tag channel is nil")
	}

	go cp.tagReader.ListenForTags(ctx)

	cp.logger.Info("Started listening for tags from reader")

Listener:
	for {
		select {
		case tagId := <-tagChannel:
			tagId = strings.ToUpper(tagId)

			return &tagId, nil
		case <-ctx.Done():
			break Listener
		default:
			cp.logger.Debug("Waiting for tag...")
			time.Sleep(time.Millisecond * 200)
		}
	}

	return nil, ctx.Err()
}

// colorMapping maps a ChargePointStatus to a color based on the indicator mapping.
func colorMapping(indicatorMapping indicator.StatusMapping, status core.ChargePointStatus) (*indicator.Color, error) {
	var color indicator.Color
	switch status {
	case core.ChargePointStatusFaulted:
		color = indicator.Color(indicatorMapping.Error)
	case core.ChargePointStatusCharging:
		color = indicator.Color(indicatorMapping.Charging)
	case core.ChargePointStatusReserved:
		color = indicator.Color(indicatorMapping.Reserved)
	case core.ChargePointStatusFinishing:
		color = indicator.Color(indicatorMapping.Finishing)
	case core.ChargePointStatusAvailable:
		color = indicator.Color(indicatorMapping.Available)
	case core.ChargePointStatusUnavailable:
		color = indicator.Color(indicatorMapping.Fault)
	default:
		return nil, errors.New("cannot find a color for the status")
	}

	return &color, nil
}

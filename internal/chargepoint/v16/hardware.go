package v16

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

// DisplayMessage sends/shows a message on the display.
func (cp *ChargePoint) DisplayMessage(message display.MessageInfo) error {
	if util.IsNilInterfaceOrPointer(cp.display) {
		cp.logger.Warn("Cannot send message to display, it is disabled or not configured")
		return nil
	}

	cp.logger.Debugf("Sending message to display: %v", message)
	go cp.display.DisplayMessage(message)
	return nil
}

func (cp *ChargePoint) indicateStatusChange(connectorIndex int, status core.ChargePointStatus) {
	logInfo := cp.logger.WithField("connector", connectorIndex+1)
	if util.IsNilInterfaceOrPointer(cp.indicator) {
		logInfo.Warn("Cannot indicate status change, indicator disabled or not configured")
		return
	}

	// Get the color for the status
	color, err := colorMapping(cp.indicatorMapping, status)
	if err != nil {
		logInfo.WithError(err).Errorf("Error indicating status")
	}

	logInfo.Debugf("Indicating connector status: %x", color)
	err = cp.indicator.ChangeColor(connectorIndex, *color)
	if err != nil {
		logInfo.WithError(err).Errorf("Error indicating status")
	}
}

// indicateCard Blinks the LED to indicate that the card was read.
func (cp *ChargePoint) indicateCard(index int, color indicator.Color) {
	logInfo := cp.logger.WithField("connector", index+1)
	if util.IsNilInterfaceOrPointer(cp.indicator) {
		logInfo.Warn("Cannot indicate card read, disabled or not configured")
		return
	}

	logInfo.Debug("Indicating a tag was read")
	err := cp.indicator.Blink(index, 3, color)
	if err != nil {
		logInfo.WithError(err).Errorf("Could not indicate a tag was read")
	}
}

// ListenForTag Listen for an RFID/NFC tag on a separate thread. If a tag is detected, call the HandleChargingRequest.
// When the tag is read, blink the LED if indication is enabled.
func (cp *ChargePoint) ListenForTag(ctx context.Context, tagChannel <-chan string) (*string, error) {
	if util.IsNilInterfaceOrPointer(tagChannel) {
		return nil, nil
	}

	if util.IsNilInterfaceOrPointer(cp.tagReader) {
		return nil, nil
	}

	go cp.tagReader.ListenForTags(ctx)

	cp.logger.Info("Started listening for tags from reader")

Listener:
	for {
		select {
		case tagId := <-tagChannel:
			tagId = strings.ToUpper(tagId)

			go func() {
				// todo _ = cp.DisplayMessage(message)
				cp.indicateCard(len(cp.evseManager.GetEVSEs()), indicator.White)
			}()

			return &tagId, nil
		case <-ctx.Done():
			break Listener
		default:
			cp.logger.Trace("Waiting for tag...")
			time.Sleep(time.Millisecond * 200)
		}
	}

	return nil, ctx.Err()
}

// colorMapping Maps a ChargePointStatus to a color based on the indicator mapping.
func colorMapping(indicatorMapping settings.IndicatorStatusMapping, status core.ChargePointStatus) (*indicator.Color, error) {
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

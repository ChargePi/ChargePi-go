package v16

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"strings"
	"time"
)

func (cp *ChargePoint) sendToLCD(messages ...string) {
	if !cp.settings.ChargePoint.Hardware.Display.IsEnabled || util.IsNilInterfaceOrPointer(cp.display) {
		return
	}

	cp.logger.Debugf("Sending message(s) to display: %v", messages)
	cp.display.DisplayMessage(chargePoint.NewMessage(time.Second*5, messages))
}

func (cp *ChargePoint) displayLEDStatus(connectorIndex int, status core.ChargePointStatus) {
	indicatorSettings := cp.settings.ChargePoint.Hardware.LedIndicator
	if !indicatorSettings.Enabled ||
		!indicatorSettings.IndicateCardRead ||
		util.IsNilInterfaceOrPointer(cp.indicator) {
		return
	}

	var color = indicator.Off
	mappings := indicatorSettings.IndicatorMappings
	mappings.Error = ""

	switch status {
	case core.ChargePointStatusFaulted:
		color = indicator.Red
	case core.ChargePointStatusCharging:
		color = indicator.Blue
	case core.ChargePointStatusReserved:
		color = indicator.Yellow
	case core.ChargePointStatusFinishing:
		color = indicator.Blue
	case core.ChargePointStatusAvailable:
		color = indicator.Green
	case core.ChargePointStatusUnavailable:
		color = indicator.Orange
	default:
		cp.logger.Error("Cannot find a color for the status")
		return
	}

	cp.logger.Debugf("Indicating connector status: %x", color)

	err := cp.indicator.DisplayColor(connectorIndex, uint32(color))
	if err != nil {
		cp.logger.WithError(err).Errorf("Error indicating status")
	}
}

// indicateCard Blinks the LED to indicate that the card was read.
func (cp *ChargePoint) indicateCard(index int, color uint32) {
	if !cp.settings.ChargePoint.Hardware.LedIndicator.Enabled || util.IsNilInterfaceOrPointer(cp.indicator) {
		return
	}

	cp.logger.Trace("Indicating tag was read")

	err := cp.indicator.Blink(index, 3, color)
	if err != nil {
		cp.logger.WithError(err).Errorf("Could not indicate card was read")
	}
}

// ListenForTag Listen for an RFID/NFC tag on a separate thread. If a tag is detected, call the HandleChargingRequest.
// Blink the LED if indication is enabled.
func (cp *ChargePoint) ListenForTag(ctx context.Context, tagChannel <-chan string) (*string, error) {
	if tagChannel == nil {
		return nil, nil
	}

	go cp.tagReader.ListenForTags(ctx)

	cp.logger.Info("Started listening for tags from reader")

Listener:
	for {
		select {
		case tagId := <-tagChannel:
			tagId = strings.ToUpper(tagId)
			go cp.indicateCard(len(cp.connectorManager.GetEVSEs()), indicator.White)
			go cp.sendToLCD("Read tag:", tagId)
			return &tagId, nil
		case <-ctx.Done():
			break Listener
		default:
			fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
			time.Sleep(time.Millisecond * 200)
		}
	}

	return nil, ctx.Err()
}

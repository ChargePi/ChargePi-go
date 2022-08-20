package v16

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/models"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"strings"
	"time"
)

func (cp *ChargePoint) sendToLCD(messages ...string) {
	if util.IsNilInterfaceOrPointer(cp.display) || !cp.settings.ChargePoint.Hardware.Display.IsEnabled {
		return
	}

	cp.logger.Debugf("Sending message(s) to display: %v", messages)
	cp.display.DisplayMessage(models.NewMessage(time.Second*5, messages))
}

func (cp *ChargePoint) displayLEDStatus(connectorIndex int, status core.ChargePointStatus) {
	if !cp.settings.ChargePoint.Hardware.LedIndicator.Enabled || util.IsNilInterfaceOrPointer(cp.indicator) {
		return
	}

	var color = indicator.Off

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
func (cp *ChargePoint) ListenForTag(ctx context.Context, tagChannel <-chan string) {
	if tagChannel == nil {
		return
	}

	cp.logger.Info("Started listening for tags from reader")

Listener:
	for {
		select {
		case tagId := <-tagChannel:
			go cp.indicateCard(len(cp.connectorManager.GetEVSEs()), indicator.White)
			go cp.sendToLCD("Read tag:", tagId)
			_, _ = cp.HandleChargingRequest(strings.ToUpper(tagId))
		case <-ctx.Done():
			break Listener
		default:
			fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
			time.Sleep(time.Millisecond * 200)
		}
	}
}

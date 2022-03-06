package v16

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/hardware/indicator"
	"strings"
	"time"
)

func (cp *ChargePoint) sendToLCD(messages ...string) {
	if cp.LCD == nil || cp.LCD.GetLcdChannel() == nil || !cp.Settings.ChargePoint.Hardware.Lcd.IsEnabled {
		return
	}

	cp.logger.Debugf("Sending message(s) to LCD: %v", messages)
	cp.LCD.GetLcdChannel() <- display.NewMessage(time.Second*5, messages)
}

func (cp *ChargePoint) displayLEDStatus(connectorIndex int, status core.ChargePointStatus) {
	if !cp.Settings.ChargePoint.Hardware.LedIndicator.Enabled || cp.Indicator == nil {
		return
	}

	var color = indicator.Off

	switch status {
	case core.ChargePointStatusFaulted:
		color = indicator.Red
		break
	case core.ChargePointStatusCharging:
		color = indicator.Blue
		break
	case core.ChargePointStatusReserved:
		color = indicator.Yellow
		break
	case core.ChargePointStatusFinishing:
		color = indicator.Blue
		break
	case core.ChargePointStatusAvailable:
		color = indicator.Green
		break
	case core.ChargePointStatusUnavailable:
		color = indicator.Orange
		break
	default:
		return
	}

	cp.logger.Debugf("Indicating connector status: %x", color)

	go func() {
		err := cp.Indicator.DisplayColor(connectorIndex, uint32(color))
		if err != nil {
			cp.logger.WithError(err).Errorf("Error indicating status")
		}
	}()
}

// indicateCard Blinks the LED to indicate that the card was read.
func (cp *ChargePoint) indicateCard(index int, color uint32) {
	if !cp.Settings.ChargePoint.Hardware.LedIndicator.Enabled || cp.Indicator == nil {
		return
	}

	cp.logger.Trace("Indicating tag was read")

	err := cp.Indicator.Blink(index, 3, color)
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
			go cp.indicateCard(len(cp.connectorManager.GetConnectors()), indicator.White)
			go cp.sendToLCD("Read tag:", tagId)
			_, _ = cp.HandleChargingRequest(strings.ToUpper(tagId))
			break
		case <-ctx.Done():
			break Listener
		default:
			fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
			time.Sleep(time.Millisecond * 200)
		}
	}
}

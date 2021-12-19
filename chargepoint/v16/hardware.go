package v16

import (
	"context"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/indicator"
	"strings"
	"time"
)

func (handler *ChargePointHandler) sendToLCD(messages ...string) {
	if handler.LCD == nil || handler.LCD.GetLcdChannel() == nil || !handler.Settings.ChargePoint.Hardware.Lcd.IsEnabled {
		return
	}

	log.Infof("Sending message(s) to LCD: %v", messages)
	handler.LCD.GetLcdChannel() <- display.NewMessage(time.Second*5, messages)
}

func (handler *ChargePointHandler) displayLEDStatus(connectorIndex int, status core.ChargePointStatus) {
	if !handler.Settings.ChargePoint.Hardware.LedIndicator.Enabled || handler.Indicator == nil {
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

	log.Infof("Indicating connector status")

	go func() {
		err := handler.Indicator.DisplayColor(connectorIndex, uint32(color))
		if err != nil {
			log.Errorf("Error indicating status: %v", err)
		}
	}()
}

// indicateCard Blinks the LED to indicate that the card was read.
func (handler *ChargePointHandler) indicateCard(index int, color uint32) {
	if !handler.Settings.ChargePoint.Hardware.LedIndicator.Enabled || handler.Indicator == nil {
		return
	}

	log.Info("Indicating tag was read")

	err := handler.Indicator.Blink(index, 3, color)
	if err != nil {
		log.Println(err)
	}
}

// ListenForTag Listen for an RFID/NFC tag on a separate thread. If a tag is detected, call the HandleChargingRequest.
// Blink the LED if indication is enabled.
func (handler *ChargePointHandler) ListenForTag(ctx context.Context) {
	if !handler.Settings.ChargePoint.Hardware.TagReader.IsEnabled || handler.TagReader == nil {
		return
	}

	log.Info("Started listening for tags from reader")

Listener:
	for {
		select {
		case tagId := <-handler.TagReader.GetTagChannel():
			go handler.indicateCard(len(handler.connectorManager.GetConnectors()), indicator.White)
			go handler.sendToLCD("Read tag:", tagId)
			handler.HandleChargingRequest(strings.ToUpper(tagId))
			break
		case <-ctx.Done():
			break Listener
		default:
			fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
			time.Sleep(time.Millisecond * 200)
		}
	}
}

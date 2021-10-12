package chargepoint

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"strings"
	"time"
)

func (handler *ChargePointHandler) sendToLCD(messages ...string) {
	if handler.LCD == nil || handler.LCD.GetLcdChannel() == nil || !handler.Settings.ChargePoint.Hardware.Lcd.IsSupported {
		return
	}

	log.Println("Sending message(s) to LCD:", messages)
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
	}

	if color != indicator.Off {
		go func() {
			err := handler.Indicator.DisplayColor(connectorIndex, uint32(color))
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

// indicateCard Blinks the LED to indicate that the card was read.
func (handler *ChargePointHandler) indicateCard(index int, color uint32) {
	if !handler.Settings.ChargePoint.Hardware.LedIndicator.Enabled || handler.Indicator == nil {
		return
	}
	_, err := scheduler.GetScheduler().Every(1).Milliseconds().LimitRunsTo(1).Do(handler.Indicator.Blink, index, 3, color)
	if err != nil {
		log.Printf("Error indicating card: %v", err)
		return
	}
}

// listenForTag Listen for an RFID/NFC tag on a separate thread. If a tag is detected, call the HandleChargingRequest.
// Blink the LED if indication is enabled.
func (handler *ChargePointHandler) listenForTag() {
	if !handler.Settings.ChargePoint.Hardware.TagReader.IsSupported {
		return
	}
	for {
		fmt.Printf("%s: Waiting for a tag \n", time.Now().String())
		select {
		case tagId := <-handler.TagReader.GetTagChannel():
			handler.indicateCard(len(handler.Connectors), indicator.White)
			handler.sendToLCD("Read tag:", tagId)
			handler.HandleChargingRequest(strings.ToUpper(tagId))
			continue
		default:
			time.Sleep(time.Millisecond * 200)
			continue
		}
	}
}

func getTypesToSample() []types2.Measurand {
	var measurands []types2.Measurand
	// Get the types to sample
	measurandsString, err := settings.GetConfigurationValue("MeterValuesSampledData")
	if err != nil {
		measurandsString = string(types2.MeasurandPowerActiveExport)
	}
	for _, s := range strings.Split(measurandsString, ",") {
		measurands = append(measurands, types2.Measurand(s))
	}
	return measurands
}

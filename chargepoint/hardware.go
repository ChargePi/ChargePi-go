package chargepoint

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"strings"
	"time"
)

func (handler *ChargePointHandler) sendToLCD(message hardware.LCDMessage) {
	if handler.LCD == nil || handler.LCD.LCDChannel == nil || !handler.Settings.ChargePoint.Hardware.Lcd.IsSupported {
		return
	}
	handler.LCD.LCDChannel <- message
}

func (handler *ChargePointHandler) displayLEDStatus(connectorIndex int, status core.ChargePointStatus) {
	if !handler.Settings.ChargePoint.Hardware.LedIndicator.Enabled || handler.LEDStrip == nil {
		return
	}
	var color = 0x00
	switch status {
	case core.ChargePointStatusFaulted:
		color = hardware.RED
		break
	case core.ChargePointStatusCharging:
		color = hardware.BLUE
		break
	case core.ChargePointStatusReserved:
		color = hardware.YELLOW
		break
	case core.ChargePointStatusFinishing:
		color = hardware.BLUE
		break
	case core.ChargePointStatusAvailable:
		color = hardware.GREEN
		break
	case core.ChargePointStatusUnavailable:
		color = hardware.ORANGE
		break
	}
	if color != 0x00 {
		go func() {
			err := handler.LEDStrip.DisplayColor(connectorIndex, uint32(color))
			if err != nil {
				log.Println(err)
				return
			}
		}()
	}
}

// indicateCard Blinks the LED to indicate that the card was read.
func (handler *ChargePointHandler) indicateCard(index int, color uint32) {
	if !handler.Settings.ChargePoint.Hardware.LedIndicator.Enabled || handler.LEDStrip == nil {
		return
	}
	_, err := scheduler.Every(1).Milliseconds().LimitRunsTo(1).Do(handler.LEDStrip.Blink, index, 3, color)
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
		case tagId := <-handler.TagReader.TagChannel:
			handler.indicateCard(len(handler.Connectors), hardware.WHITE)
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

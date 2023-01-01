package v16

import (
	"context"
	"fmt"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"strings"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/indicator"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

func (cp *ChargePoint) sendToLCD(messages ...string) {
	if util.IsNilInterfaceOrPointer(cp.display) {
		cp.logger.Warn("Cannot send message to display, it is disabled or not configured")
		return
	}

	cp.logger.Debugf("Sending message(s) to display: %v", messages)
	go cp.display.DisplayMessage(chargePoint.NewMessage(time.Second*5, messages))
}

func (cp *ChargePoint) displayStatusChangeOnIndicator(connectorIndex int, status core.ChargePointStatus) {
	logInfo := cp.logger.WithField("connector", connectorIndex+1)
	if util.IsNilInterfaceOrPointer(cp.indicator) {
		logInfo.Warn("Cannot indicate status change, indicator disabled or not configured")
		return
	}

	var color = indicator.Off
	switch status {
	case core.ChargePointStatusFaulted:
		color = indicator.Color(cp.indicatorMapping.Error)
	case core.ChargePointStatusCharging:
		color = indicator.Color(cp.indicatorMapping.Charging)
	case core.ChargePointStatusReserved:
		color = indicator.Color(cp.indicatorMapping.Reserved)
	case core.ChargePointStatusFinishing:
		color = indicator.Color(cp.indicatorMapping.Finishing)
	case core.ChargePointStatusAvailable:
		color = indicator.Color(cp.indicatorMapping.Available)
	case core.ChargePointStatusUnavailable:
		color = indicator.Color(cp.indicatorMapping.Fault)
	default:
		logInfo.Error("Cannot find a color for the status")
		return
	}

	logInfo.Debugf("Indicating connector status: %x", color)
	err := cp.indicator.DisplayColor(connectorIndex, color)
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

			go func() {
				cp.sendToLCD("Read tag:", tagId)
				cp.indicateCard(len(cp.connectorManager.GetEVSEs()), indicator.White)
			}()

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

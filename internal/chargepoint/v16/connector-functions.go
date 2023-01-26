package v16

import (
	"context"
	"time"

	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display/i18n"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
)

// restoreState After connecting to the central system, try to restore the previous state of each EVSE and notify
// the system about its state.
//
// If the ConnectorStatus was "Preparing" or "Charging", try to resume or start charging.
// If the charging fails, change the connector status and notify the central system.
func (cp *ChargePoint) restoreState() {
	cp.logger.Info("Restoring evses' state")
	err := cp.evseManager.RestoreEVSEs()
	if err != nil {
		cp.logger.WithError(err).Error("Unable to restore states")
	}
}

// notifyConnectorStatus Notify the central system about the connector's status and updates the LED indicator.
func (cp *ChargePoint) notifyConnectorStatus(evseId int, status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	request := core.NewStatusNotificationRequest(evseId, errCode, status)
	request.Timestamp = types.NewDateTime(time.Now())

	callback := func(confirmation ocpp.Response, protoError error) {
		cp.logger.WithField("evseId", evseId).Infof("Notified status %s of the connector", status)
	}

	err := util.SendRequest(cp.chargePoint, request, callback)
	util.HandleRequestErr(err, "Cannot send status of connector")
}

// ListenForConnectorStatusChange listen for change in connector and notify the central system about the state
func (cp *ChargePoint) ListenForConnectorStatusChange(ctx context.Context, ch <-chan notifications.StatusNotification) {
	cp.logger.Debug("Starting to listen for connector status change")

Listener:
	for {
		select {
		case c := <-ch:
			// Connector starts with index 1
			connectorIndex := c.EvseId - 1
			status := core.ChargePointStatus(c.Status)
			errCode := core.ChargePointErrorCode(c.Status)

			logInfo := cp.logger.WithFields(log.Fields{
				"evseId":  c.EvseId,
				"status":  status,
				"errCode": errCode,
			})
			logInfo.Info("Received evse status update")

			go func() {
				logInfo.Info("Displaying status change on hardware")
				// Display connector status change on display and indicator
				cp.displayStatusChangeOnIndicator(connectorIndex, status)
				cp.displayStatusChangeOnDisplay(c.EvseId, status)
			}()

			go cp.handleStatusUpdate(ctx, c.EvseId, status)

			// Send a status notification to the Central System
			cp.notifyConnectorStatus(c.EvseId, status, errCode)

		case meterVal := <-cp.meterValuesChannel:
			logInfo := cp.logger.WithFields(log.Fields{
				"evseId":      meterVal.EvseId,
				"transaction": meterVal.TransactionId,
			})
			logInfo.Info("Received meter value update")

			// Send a meter value notification to the Central System
			values := core.NewMeterValuesRequest(meterVal.EvseId, meterVal.MeterValues)
			err := util.SendRequest(cp.chargePoint, values, func(confirmation ocpp.Response, protoError error) {
				if protoError != nil {
					return
				}

				logInfo.Info("Sent a meter value update")
			})
			if err != nil {
				logInfo.WithError(err).Errorf("Cannot send meter values")
			}

		case <-ctx.Done():
			break Listener
		}
	}
}

func (cp *ChargePoint) displayStatusChangeOnDisplay(connectorId int, status core.ChargePointStatus) {
	var (
		language = ""
		message  []string
		err      error
	)

	switch cp.display.GetType() {
	case display.DriverHD44780:

		switch status {
		case core.ChargePointStatusAvailable:
			message, err = i18n.TranslateConnectorAvailableMessage(language, connectorId)
		case core.ChargePointStatusFinishing:
			message, err = i18n.TranslateConnectorFinishingMessage(language, connectorId)
		case core.ChargePointStatusCharging:
			message, err = i18n.TranslateConnectorChargingMessage(language, connectorId)
		case core.ChargePointStatusFaulted:
			message, err = i18n.TranslateConnectorFaultedMessage(language, connectorId)
		default:
			err = display.ErrDisplayUnsupported
		}

	default:
		err = display.ErrDisplayUnsupported
	}

	if err != nil {
		cp.logger.WithError(err).Errorf("Error displaying status")
		return
	}

	cp.sendToLCD(message...)
}

// handleStatusUpdate if an EV is connected, ask for authentication, if it was disconnected, stop the transaction.
func (cp *ChargePoint) handleStatusUpdate(ctx context.Context, evseId int, status core.ChargePointStatus) {

	// Todo design a plugin interface so it can be called here

	switch status {
	case core.ChargePointStatusPreparing:

		// Listen for a tag for a minute. If the tag is tapped, start charging.
		listenCtx, cancel := context.WithTimeout(ctx, time.Minute)
		tag, err := cp.ListenForTag(listenCtx, cp.tagReader.GetTagChannel())
		if err == nil {
			// Tag was found, start charging
			err = cp.StartCharging(evseId, 1, *tag)
			if err != nil {
				cp.logger.WithError(err).Error("Cannot start charging")
			}
		}

		cancel()
	case core.ChargePointStatusSuspendedEV, core.ChargePointStatusSuspendedEVSE:
		stopTransactionOnEVDisconnect, err := ocppConfigManager.GetConfigurationValue(configuration.StopTransactionOnEVSideDisconnect.String())

		// If EV disconnects and the StopTransactionOnEVDisconnect is enabled, stop the transaction
		if stopTransactionOnEVDisconnect != nil && *stopTransactionOnEVDisconnect == "true" {
			stopChargingErr := cp.StopCharging(evseId, 1, core.ReasonEVDisconnected)
			if stopChargingErr != nil {
				cp.logger.WithError(err).Error("Cannot stop charging")
			}
		}

	case core.ChargePointStatusFaulted:
		err := cp.StopCharging(evseId, 1, core.ReasonEmergencyStop)
		if err != nil {
			cp.logger.WithError(err).Error("Cannot stop charging")
		}
	}
}

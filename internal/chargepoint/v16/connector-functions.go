package v16

import (
	"context"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/hardware/display/i18n"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
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
		cp.logger.WithError(err).Warn("Unable to restore states")
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
				logInfo.Debug("Displaying status change")
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
				logInfo.Info("Sent a meter value update")
			})
			if err != nil {
				logInfo.WithError(err).Errorf("Cannot send meter values")
			}

			// todo add plugin/middleware support

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
		case core.ChargePointStatusUnavailable:
		}

	default:
		err = display.ErrDisplayUnsupported
	}

	if err != nil {
		cp.logger.WithError(err).Errorf("Error displaying status")
		return
	}

	_ = cp.DisplayMessage(notifications.NewMessage(time.Second*10, message...))
}

// handleStatusUpdate if an EV is connected, ask for authentication, if it was disconnected, stop the transaction.
func (cp *ChargePoint) handleStatusUpdate(ctx context.Context, evseId int, status core.ChargePointStatus) {
	logInfo := cp.logger.WithField("evseId", evseId).WithField("status", status)
	logInfo.Debug("Handling status update")
	// Todo design a plugin/middleware interface so it can be called here

	switch status {
	case core.ChargePointStatusAvailable:
		// Todo indicate that the EVSE is available

	case core.ChargePointStatusPreparing:

		// Check if FreeMode is enabled. This will bypass any authentication requirement.
		if cp.info.FreeMode {
			cp.startChargingFreeMode(evseId)
			break
		}

		// Default to RFID authentication
		cp.authenticateWithRfidCard(ctx, evseId)

	case core.ChargePointStatusCharging:
		// Todo indicate that the EV is charging
	case core.ChargePointStatusSuspendedEV:
		// Todo indicate that the EV has halted charging
	case core.ChargePointStatusSuspendedEVSE:
		// Todo indicate that the EVSE has halted charging
	case core.ChargePointStatusFinishing:
		stopTransactionOnEVDisconnect, err := ocppConfigManager.GetConfigurationValue(configuration.StopTransactionOnEVSideDisconnect.String())

		// If EV disconnects and the StopTransactionOnEVDisconnect is enabled, stop the transaction
		if stopTransactionOnEVDisconnect != nil && *stopTransactionOnEVDisconnect == "true" {
			stopChargingErr := cp.StopCharging(evseId, 1, core.ReasonEVDisconnected)
			if stopChargingErr != nil {
				logInfo.WithError(err).Error("Cannot stop charging")
				// Todo Indicate that the charging hasn't been successfully stopped
			}

			// Todo Indicate that the charging has been stopped
		}

	case core.ChargePointStatusFaulted:
		err := cp.StopCharging(evseId, 1, core.ReasonEmergencyStop)
		if err != nil {
			logInfo.WithError(err).Error("Cannot stop charging")
		}
	}
}

func (cp *ChargePoint) startChargingFreeMode(evseId int) {
	logInfo := cp.logger.WithField("evseId", evseId)
	logInfo.Info("Free mode enabled, starting charging")

	err := cp.evseManager.StartCharging(evseId, nil)
	if err != nil {
		logInfo.WithError(err).Errorf("Unable to start charging connector")
	}
}

func (cp *ChargePoint) authenticateWithRfidCard(ctx context.Context, evseId int) {
	logInfo := cp.logger.WithField("evseId", evseId)

	// Listen for a tag for a minute. If the tag is tapped, start charging.
	listenCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	cp.logger.Info("Waiting for the user to tap a tag")

	tag, err := cp.ListenForTag(listenCtx, cp.tagReader.GetTagChannel())
	switch err {
	case nil:
		// Tag was found, start charging
		err = cp.StartCharging(evseId, 1, *tag)
		if err != nil {
			logInfo.WithError(err).Error("Cannot start charging")
		}
		// Indicate charging should start

	case context.DeadlineExceeded:
		// Indicate timeout
	default:
		logInfo.WithError(err).Error("Error while listening for tag")
		// Indicate error
	}
}

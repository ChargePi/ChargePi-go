package v16

import (
	"context"
	"errors"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/notifications"
	"github.com/xBlaz3kx/ChargePi-go/pkg/display"
	"github.com/xBlaz3kx/ocppManager-go/ocpp_v16"
)

// restoreState After connecting to the central system, try to restore the previous state of each EVSE and notify
// the system about its state.
//
// If the ConnectorStatus was "Preparing" or "Charging", try to resume or start charging.
// If the charging fails, change the connector status and notify the central system.
func (cp *ChargePoint) restoreState() {
	cp.logger.Info("Restoring charge point state")
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

	err := cp.sendRequest(request, callback)
	cp.handleRequestErr(err, "Cannot send status of connector")
}

// ListenForConnectorStatusChange listen for change in connector and notify the central system about the state
func (cp *ChargePoint) ListenForConnectorStatusChange(ctx context.Context, ch <-chan notifications.StatusNotification) {
	cp.logger.Debug("Starting to listen for connector status change")

Listener:
	for {
		select {
		case c := <-ch:
			// Connector starts with index 1
			status := core.ChargePointStatus(c.Status)
			errCode := core.ChargePointErrorCode(c.Status)

			logInfo := cp.logger.WithFields(log.Fields{
				"evseId":  c.EvseId,
				"status":  status,
				"errCode": errCode,
			})
			logInfo.Info("Received evse status update")

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
			err := cp.sendRequest(values, func(confirmation ocpp.Response, protoError error) {
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

// displayStatusChangeOnDisplay sends an update to the Display based on the EVSE status change.
func (cp *ChargePoint) displayStatusChangeOnDisplay(connectorId int, status core.ChargePointStatus) {

	switch cp.display.GetType() {
	case display.DriverHD44780:

	default:
		cp.logger.Errorf("Display unsupported")
	}
}

// handleStatusUpdate if an EV is connected, ask for authentication, if it was disconnected, stop the transaction.
func (cp *ChargePoint) handleStatusUpdate(ctx context.Context, evseId int, status core.ChargePointStatus) {
	logInfo := cp.logger.WithField("evseId", evseId).WithField("status", status)
	logInfo.Debug("Handling status update")

	cp.indicateStatusChange(evseId-1, status)
	cp.displayStatusChangeOnDisplay(evseId-1, status)
	// Todo design a plugin/middleware interface

	switch status {
	case core.ChargePointStatusAvailable:
		// Todo indicate that the EVSE is available

	case core.ChargePointStatusPreparing:

		// Check if FreeMode is enabled. This will bypass any authentication requirement.
		if cp.info.FreeMode {
			cp.StartChargingFreeMode(evseId)
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
		// If EV disconnects and the StopTransactionOnEVDisconnect is enabled, stop the transaction
		stopTransactionOnEVDisconnect, err := cp.settingsManager.GetOcppV16Manager().GetConfigurationValue(ocpp_v16.StopTransactionOnEVSideDisconnect)
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

func (cp *ChargePoint) StartChargingFreeMode(evseId int) error {
	if !cp.info.FreeMode {
		return errors.New("free mode is not enabled")
	}

	logInfo := cp.logger.WithField("evseId", evseId)
	logInfo.Info("Free mode enabled, starting charging")

	measurements, sampleInterval := cp.getSessionParameters()

	return cp.evseManager.StartCharging(evseId, nil, measurements, sampleInterval)
}

// authenticateWithRfidCard waits for a tag to be tapped and starts charging if the tag is valid.
func (cp *ChargePoint) authenticateWithRfidCard(ctx context.Context, evseId int) {
	// Listen for a tag for a minute. If the tag is tapped, trigger the standard charging flow.
	listenCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	logInfo := cp.logger.WithField("evseId", evseId)
	logInfo.Info("Waiting for the user to tap a tag")

	tag, err := cp.ListenForTag(listenCtx, cp.tagReader.GetTagChannel())
	switch {
	case err == nil:
		// Tag was found, attempt to start charging
		err = cp.StartCharging(evseId, 1, *tag)
		if err != nil {
			logInfo.WithError(err).Error("Cannot start charging")
		}

		// Indicate charging should start
	case errors.Is(err, context.DeadlineExceeded):
		// Indicate timeout
	default:
		logInfo.WithError(err).Error("Error while listening for tag")
		// Indicate error
	}
}

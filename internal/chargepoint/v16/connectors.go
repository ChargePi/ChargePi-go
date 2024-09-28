package v16

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"time"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/indicator"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

func (cp *ChargePoint) OnUnlockConnector(request *core.UnlockConnectorRequest) (confirmation *core.UnlockConnectorConfirmation, err error) {
	cp.logger.Sugar().Infof("Received request %s", request.GetFeatureName())
	response := core.UnlockStatusNotSupported

	conn, fErr := cp.evseManager.GetEVSE(request.ConnectorId)
	switch fErr {
	case nil:
		cp.logger.Sugar().Infof("Unlocking connector %d", request.ConnectorId)
		conn.GetEvcc().Unlock()
		response = core.UnlockStatusUnlocked
	}

	return core.NewUnlockConnectorConfirmation(response), nil
}

// notifyStatus Notify the central system about the connector's status and updates the LED indicator.
func (cp *ChargePoint) notifyStatus(evseId int, status core.ChargePointStatus, errCode core.ChargePointErrorCode) {
	request := core.NewStatusNotificationRequest(evseId, errCode, status)
	request.Timestamp = types.NewDateTime(time.Now())

	callback := func(confirmation ocpp.Response, protoError error) {
		cp.logger.With(zap.Int("evse_id", evseId), zap.String("status",string(status))).Info("Notified status of the connector")
	}

	err := cp.sendRequest(request, callback)
	cp.handleRequestErr(err, "Cannot send status of connector")
}

// ListenForConnectorStatusChange listen for change in connector and notify the central system about the state
func (cp *ChargePoint) ListenForConnectorStatusChange(ctx context.Context) {
	cp.logger.Debug("Starting to listen for connector status change")
	ch := cp.evseManager.GetNotificationChannel()

Listener:
	for {
		select {
		case c := <-ch:
			// Connector starts with index 1
			status := core.ChargePointStatus(c.Status)
			errCode := core.ChargePointErrorCode(c.Status)

			logInfo := cp.logger.With(
				zap.Int("evse_id", c.EvseId),
				zap.String("status", string(status)),
				zap.String("error_code", string(errCode)))
			logInfo.Info("Received evse status update")

			// Display the status change on the charge point
			go cp.handleStatusUpdate(ctx, c.EvseId, status)

			// Send a status notification to the Central System
			cp.notifyStatus(c.EvseId, status, errCode)

		case meterVal := <-cp.meterValuesChannel:
			logInfo := cp.logger.With(
				zap.Int("evse_id", meterVal.EvseId))
			logInfo = logInfo.With(zap.Int("meter_values", len(meterVal.MeterValues)))
			logInfo.Info("Received meter value update")

			// Send a meter value notification to the Central System
			values := core.NewMeterValuesRequest(meterVal.EvseId, meterVal.MeterValues)
			err := cp.sendRequest(values, func(confirmation ocpp.Response, protoError error) {
				logInfo.Info("Sent a meter value update")
			})
			if err != nil {
				logInfo.With(zap.Error(err)).Error("Cannot send meter values")
			}
		case <-ctx.Done():
			// todo wait for all the messages in other channels to be processed?
			break Listener
		}
	}
}

// displayStatusChangeOnDisplay sends an update to the Display based on the EVSE status change.
func (cp *ChargePoint) displayStatusChangeOnDisplay(connectorId int, status core.ChargePointStatus) {
	logInfo := cp.logger.With(zap.Int("connector", connectorId+1), zap.String("status", string(status)))
	logInfo.Debug("Updating status on display")
	// todo
}

// handleStatusUpdate if an EV is connected, ask for authentication, if it was disconnected, stop the transaction.
func (cp *ChargePoint) handleStatusUpdate(ctx context.Context, evseId int, status core.ChargePointStatus) {
	logger := cp.logger.With(zap.Int("evse_id", evseId), zap.String("status", string(status)))
	logger.Debug("Handling status update")

	cp.indicateStatusChange(evseId-1, status)
	cp.displayStatusChangeOnDisplay(evseId-1, status)

	switch status {
	case core.ChargePointStatusAvailable:
		// Todo indicate that the EVSE is available

	case core.ChargePointStatusPreparing:

		// Check if FreeMode is enabled. This will bypass any authentication requirement.
		if cp.info.FreeMode {
			err := cp.StartChargingFreeMode(evseId)
			if err != nil {
				logger.With(zap.Error(err)).Error("Cannot start charging in free mode")
			}

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
		stopTransactionOnEVDisconnect, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.StopTransactionOnEVSideDisconnect)
		if stopTransactionOnEVDisconnect != nil && *stopTransactionOnEVDisconnect == "true" {
			stopChargingErr := cp.StopCharging(evseId, 1, core.ReasonEVDisconnected)
			if stopChargingErr != nil {
				logger.With(zap.Error(err)).Error("Cannot stop charging")
				// Todo Indicate that the charging hasn't been successfully stopped
			}

			// Todo Indicate that the charging has been stopped
		}

	case core.ChargePointStatusFaulted:
		err := cp.StopCharging(evseId, 1, core.ReasonEmergencyStop)
		if err != nil {
			logger.With(zap.Error(err)).Error("Cannot stop charging")
		}
	}
}

// authenticateWithRfidCard waits for a tag to be tapped and starts charging if the tag is valid.
func (cp *ChargePoint) authenticateWithRfidCard(ctx context.Context, evseId int) {
	// Listen for a tag for a minute. If the tag is tapped, trigger the standard charging flow.
	listenCtx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	logger := cp.logger.With(zap.Int("evse_id", evseId))
	logger.Info("Waiting for the user to tap a tag")

	tag, err := cp.ListenForTag(listenCtx)
	switch {
	case err == nil:
		// Indicate tag was read
		cp.indicateCardRead(evseId-1, indicator.Green)

		// Attempt to start charging
		err = cp.StartCharging(evseId, 1, *tag)
		if err != nil {
			logger.With(zap.Error(err)).Error("Cannot start charging")
		}
	case errors.Is(err, context.DeadlineExceeded):
		// Indicate timeout
		logger.Warn("Timeout while waiting for tag")
	default:
		logger.With(zap.Error(err)).Error("Error while listening for tag")
	}
}

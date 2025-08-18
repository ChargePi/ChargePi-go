package v16

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"

	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/pkg/util"
)

func (cp *ChargePoint) OnRemoteStopTransaction(request *core.RemoteStopTransactionRequest) (confirmation *core.RemoteStopTransactionConfirmation, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	logger := cp.logger.With(zap.Int("transactionId", request.TransactionId))
	logger.Sugar().Infof("Received request %s", request.GetFeatureName())

	response := types.RemoteStartStopStatusRejected
	transactionId := fmt.Sprintf("%d", request.TransactionId)

	session, fErr := cp.sessionService.GetSessionWithTransactionId(ctx, transactionId)
	if fErr == nil {
		cp.logger.Info("Stopping transaction")
		response = types.RemoteStartStopStatusAccepted

		// Delay stopping the transaction by 3 seconds
		_, schedulerErr := cp.scheduler.Every(3).Seconds().LimitRunsTo(1).Do(cp.StopCharging, context.Background(), session.EvseId, session.ConnectorId, core.ReasonRemote)
		if schedulerErr != nil {
			cp.logger.With(zap.Error(err)).Error("Failed to schedule stop charging")
			response = types.RemoteStartStopStatusRejected
		}
	}

	return core.NewRemoteStopTransactionConfirmation(response), nil
}

// StopCharging Stops a transaction on the specified EVSE and connector. The reason for stopping the transaction should be provided.
func (cp *ChargePoint) StopCharging(ctx context.Context, evseId, connectorId int, reason core.Reason) error {
	cpEvse, err := cp.evseManager.GetEVSE(evseId)
	if err != nil {
		return err
	}

	return cp.stopChargingConnector(ctx, cpEvse, reason)
}

// stopChargingConnector Stop charging a connector with the specified ID.
func (cp *ChargePoint) stopChargingConnector(ctx context.Context, connector evse.EVSE, reason core.Reason) error {
	if util.IsNilInterfaceOrPointer(connector) {
		return chargepoint.ErrConnectorNil
	}

	evseId := connector.GetEvseId()
	logger := cp.logger.With(zap.Int("evseId", evseId), zap.String("reason", string(reason)))
	logger.Info("Stopping transaction")

	// Check if the connector is already stopped
	session, err := cp.sessionService.GetSession(ctx, evseId, nil)
	if err != nil {
		return err
	}

	// Stop charging on EVSE
	err = connector.StopCharging(reason)
	if err != nil {
		logger.With(zap.Error(err)).Error("Unable to stop charging")
		return err
	}

	logger.Sugar().Infof("Stopped charging at %s", time.Now())

	// Mark the session as stopped in the session manager
	err = cp.sessionService.StopSession(ctx, evseId, nil, nil, &session.TransactionId)
	if err != nil {
		logger.With(zap.Error(err)).Warn("Unable to stop session")
	}

	// Notify the backend that the session has stopped
	transactionId, convErr := strconv.Atoi(session.TransactionId)
	if convErr != nil {
		return convErr
	}

	// Get the energy consumption of the session
	consumption, err := session.GetEnergyConsumption()
	if err != nil {
		return err
	}

	request := core.NewStopTransactionRequest(
		int(consumption),
		types.NewDateTime(time.Now()),
		transactionId,
	)
	request.Reason = reason

	var callback = func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logger.With(zap.Error(err)).Error("Server responded with error for stopping a transaction")
			return
		}
	}

	return cp.sendRequest(request, callback)
}

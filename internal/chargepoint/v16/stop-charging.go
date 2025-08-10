package v16

import (
	"go.uber.org/zap"
	"strconv"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/ChargePi/ChargePi-go/internal/pkg/util"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

func (cp *ChargePoint) StopCharging(evseId, connectorId int, reason core.Reason) error {
	cpEvse, err := cp.evseManager.GetEVSE(evseId)
	if err != nil {
		return err
	}

	return cp.stopChargingConnector(cpEvse, reason)
}

// stopChargingConnector Stop charging a connector with the specified ID. Update the status(es), turn off the ConnectorImpl and calculate the energy consumed.
func (cp *ChargePoint) stopChargingConnector(connector evse.EVSE, reason core.Reason) error {
	if util.IsNilInterfaceOrPointer(connector) {
		return chargePoint.ErrConnectorNil
	}

	logger := cp.logger.With()

	// Check if the connector is already stopped
	session, err := cp.sessionManager.GetSession(connector.GetEvseId(), nil)
	if err != nil {
		return err
	}

	transactionId, convErr := strconv.Atoi(session.TransactionId)
	if convErr != nil {
		return convErr
	}

	request := core.NewStopTransactionRequest(
		int(session.CalculateEnergyConsumptionWithAvgPower()),
		types.NewDateTime(time.Now()),
		transactionId,
	)
	request.Reason = reason

	var callback = func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logger.With(zap.Error(err)).Error("Server responded with error for stopping a transaction")
			return
		}

		logger.Info("Stopping transaction")

		// Stop charging on EVSE
		err = connector.StopCharging(reason)
		if err != nil {
			logger.With(zap.Error(err)).Error("Unable to stop charging")
			return
		}

		err = cp.sessionManager.StopSession(session.TransactionId)
		if err != nil {
			logger.With(zap.Error(err)).Warn("Unable to stop session")
		}

		logger.Sugar().Infof("Stopped charging at %s", time.Now())
	}

	return cp.sendRequest(request, callback)
}

package v16

import (
	"strconv"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
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

	logInfo := cp.logger.WithFields(log.Fields{
		"evseId": connector.GetEvseId(),
		"reason": reason,
	})

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
			logInfo.WithError(protoError).Errorf("Server responded with error for stopping a transaction")
			return
		}

		logInfo.Info("Stopping transaction")

		// Stop charging on EVSE
		err = connector.StopCharging(reason)
		if err != nil {
			logInfo.WithError(err).Errorf("Unable to stop charging")
			return
		}

		err = cp.sessionManager.StopSession(session.TransactionId)
		if err != nil {
			logInfo.WithError(err).Warnf("Unable to stop session")
		}

		logInfo.Infof("Stopped charging at %s", time.Now())
	}

	return util.SendRequest(cp.chargePoint, request, callback)
}

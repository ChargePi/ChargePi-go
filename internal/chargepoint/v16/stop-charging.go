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
	cpEvse, err := cp.evseManager.FindEVSE(evseId)
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

	var (
		transactionId, convErr = strconv.Atoi(connector.GetTransactionId())
		logInfo                = cp.logger.WithFields(log.Fields{
			"evseId": connector.GetEvseId(),
			"reason": reason,
		})
	)

	if convErr != nil {
		return convErr
	}

	session := connector.GetSession()
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
		err := connector.StopCharging(reason)
		if err != nil {
			logInfo.WithError(err).Errorf("Unable to stop charging")
			return
		}

		logInfo.Infof("Stopped charging at %s", time.Now())
	}

	return util.SendRequest(cp.chargePoint, request, callback)
}

// stopChargingConnectorWithTagId Search for a connector with the tagId and stop the charging.
func (cp *ChargePoint) stopChargingConnectorWithTagId(tagId string, reason core.Reason) error {
	var c, err = cp.evseManager.FindEVSEWithTagId(tagId)
	if err != nil {
		return err
	}

	return cp.stopChargingConnector(c, reason)
}

// stopChargingConnectorWithTransactionId Search for a connector with the transactionId and stop the charging.
func (cp *ChargePoint) stopChargingConnectorWithTransactionId(transactionId string) error {
	var c, err = cp.evseManager.FindEVSEWithTransactionId(transactionId)
	if err != nil {
		return err
	}

	return cp.stopChargingConnector(c, core.ReasonRemote)
}

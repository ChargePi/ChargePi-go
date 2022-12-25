package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint/components/evse"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"strconv"
	"time"
)

func (cp *ChargePoint) StopCharging(evseId, connectorId int, reason core.Reason) error {
	return nil
}

// stopChargingConnector Stop charging a connector with the specified ID. Update the status(es), turn off the ConnectorImpl and calculate the energy consumed.
func (cp *ChargePoint) stopChargingConnector(connector evse.EVSE, reason core.Reason) error {
	if util.IsNilInterfaceOrPointer(connector) {
		return chargePoint.ErrConnectorNil
	}

	var (
		stopTransactionOnEVDisconnect, err = ocppConfigManager.GetConfigurationValue(configuration.StopTransactionOnEVSideDisconnect.String())
		transactionId, convErr             = strconv.Atoi(connector.GetTransactionId())
		logInfo                            = cp.logger.WithFields(log.Fields{
			"evseId": connector.GetEvseId(),
			"reason": reason,
		})
	)

	if convErr != nil {
		return convErr
	}

	if stopTransactionOnEVDisconnect != nil && *stopTransactionOnEVDisconnect != "true" && reason == core.ReasonEVDisconnected {
		logInfo.Info("The charging without stopping the transaction")
		return connector.StopCharging(reason)
	}

	request := core.NewStopTransactionRequest(
		int(connector.CalculateSessionAvgEnergyConsumption()),
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
		err = connector.StopCharging(reason)
		if err != nil {
			logInfo.WithError(err).Errorf("Unable to stop charging")
			return
		}

		logInfo.Infof("Stopped charging at %s", time.Now())
	}

	return util.SendRequest(cp.chargePoint, request, callback)
}

// stopChargingConnectorWithTagId Search for a ConnectorImpl that contains the tagId and stop the charging.
func (cp *ChargePoint) stopChargingConnectorWithTagId(tagId string, reason core.Reason) error {
	var c, err = cp.connectorManager.FindEVSEWithTagId(tagId)
	if err != nil {
		return err
	}

	return cp.stopChargingConnector(c, reason)
}

// stopChargingConnectorWithTransactionId Search for a ConnectorImpl that contains the transactionId and stop the charging.
func (cp *ChargePoint) stopChargingConnectorWithTransactionId(transactionId string) error {
	var c, err = cp.connectorManager.FindEVSEWithTransactionId(transactionId)
	if err != nil {
		return err
	}

	return cp.stopChargingConnector(c, core.ReasonRemote)
}

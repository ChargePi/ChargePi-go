package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/errors"
	"github.com/xBlaz3kx/ChargePi-go/pkg/util"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"strconv"
	"time"
)

// stopChargingConnector Stop charging a connector with the specified ID. Update the status(es), turn off the ConnectorImpl and calculate the energy consumed.
func (cp *ChargePoint) stopChargingConnector(connector connector.Connector, reason core.Reason) error {
	if util.IsNilInterfaceOrPointer(connector) {
		return errors.ErrConnectorNil
	}

	var (
		stopTransactionOnEVDisconnect, err = ocppConfigManager.GetConfigurationValue(v16.StopTransactionOnEVSideDisconnect.String())
		transactionId, convErr             = strconv.Atoi(connector.GetTransactionId())
		logInfo                            = log.WithFields(log.Fields{
			"evseId":      connector.GetEvseId(),
			"connectorId": connector.GetConnectorId(),
			"reason":      reason,
		})
	)

	if err != nil {
		stopTransactionOnEVDisconnect = "true"
	}

	if convErr != nil {
		return convErr
	}

	if !(connector.IsCharging() || connector.IsPreparing()) {
		return errors.ErrConnectorNotCharging
	}

	if stopTransactionOnEVDisconnect != "true" && reason == core.ReasonEVDisconnected {
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

		schedulerErr := cp.scheduler.RemoveByTag(fmt.Sprintf("connector%dSampling", connector.GetConnectorId()))
		if schedulerErr != nil {
			logInfo.WithError(err).Errorf("Cannot remove sampling schedule")
		}

		schedulerErr = cp.scheduler.RemoveByTag(fmt.Sprintf("connector%dTimer", connector.GetConnectorId()))
		if schedulerErr != nil {
			logInfo.WithError(err).Errorf("Cannot remove stop charging schedule")
		}

		logInfo.Infof("Stopped charging at %s", time.Now())
	}

	return util.SendRequest(cp.chargePoint, request, callback)
}

// stopChargingConnectorWithTagId Search for a ConnectorImpl that contains the tagId and stop the charging.
func (cp *ChargePoint) stopChargingConnectorWithTagId(tagId string, reason core.Reason) error {
	var c = cp.connectorManager.FindConnectorWithTagId(tagId)
	if !util.IsNilInterfaceOrPointer(c) {
		return cp.stopChargingConnector(c, reason)
	}

	return errors.ErrNoConnectorWithTag
}

// stopChargingConnectorWithTransactionId Search for a ConnectorImpl that contains the transactionId and stop the charging.
func (cp *ChargePoint) stopChargingConnectorWithTransactionId(transactionId string) error {
	var c = cp.connectorManager.FindConnectorWithTransactionId(transactionId)
	if !util.IsNilInterfaceOrPointer(c) {
		return cp.stopChargingConnector(c, core.ReasonRemote)
	}

	return errors.ErrNoConnectorWithTransaction
}

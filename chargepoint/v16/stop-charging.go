package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/conf-manager"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"strconv"
	"time"
)

// stopChargingConnector Stop charging a connector with the specified ID. Update the status(es), turn off the ConnectorImpl and calculate the energy consumed.
func (handler *ChargePointHandler) stopChargingConnector(connector connector.Connector, reason core.Reason) error {
	if data.IsNilInterfaceOrPointer(connector) {
		return chargepoint.ErrConnectorNil
	}

	var (
		stopTransactionOnEVDisconnect, err = conf_manager.GetConfigurationValue("StopTransactionOnEVSideDisconnect")
		transactionId, convErr             = strconv.Atoi(connector.GetTransactionId())
		logInfo                            = log.WithFields(log.Fields{
			"evseId":      connector.GetEvseId(),
			"connectorId": connector.GetConnectorId(),
			"reason":      reason,
		})
	)

	if err != nil {
		return err
	}

	if convErr != nil {
		return convErr
	}

	if connector.IsCharging() || connector.IsPreparing() {
		if stopTransactionOnEVDisconnect != "true" && reason == core.ReasonEVDisconnected {
			return connector.StopCharging(reason)
		}

		request := core.NewStopTransactionRequest(
			int(connector.CalculateSessionAvgEnergyConsumption()),
			types2.NewDateTime(time.Now()),
			transactionId,
		)
		request.Reason = reason

		var callback = func(confirmation ocpp.Response, protoError error) {
			if protoError != nil {
				logInfo.Infof("Server responded with error for stopping a transaction: %v", err)
				return
			}

			logInfo.Info("Stopping transaction")
			err = connector.StopCharging(reason)
			if err != nil {
				logInfo.Errorf("Unable to stop charging connector: %v", err)
				return
			}

			schedulerErr := scheduler.GetScheduler().RemoveByTag(fmt.Sprintf("connector%dSampling", connector.GetConnectorId()))
			if schedulerErr != nil {
				log.Errorf("%v", err)
			}

			schedulerErr = scheduler.GetScheduler().RemoveByTag(fmt.Sprintf("connector%dTimer", connector.GetConnectorId()))
			if schedulerErr != nil {
				log.Errorf("%v", err)
			}

			logInfo.Infof("Stopped charging connector at %s", time.Now())
		}

		return handler.SendRequest(request, callback)
	}

	return chargepoint.ErrConnectorNotCharging
}

// stopChargingConnectorWithTagId Search for a ConnectorImpl that contains the tagId and stop the charging.
func (handler *ChargePointHandler) stopChargingConnectorWithTagId(tagId string, reason core.Reason) error {
	var c = handler.connectorManager.FindConnectorWithTagId(tagId)
	if c != nil {
		return handler.stopChargingConnector(c, reason)
	}

	return chargepoint.ErrNoConnectorWithTag
}

// stopChargingConnectorWithTransactionId Search for a ConnectorImpl that contains the transactionId and stop the charging.
func (handler *ChargePointHandler) stopChargingConnectorWithTransactionId(transactionId string) error {
	var c = handler.connectorManager.FindConnectorWithTransactionId(transactionId)
	if c != nil {
		return handler.stopChargingConnector(c, core.ReasonRemote)
	}

	return chargepoint.ErrNoConnectorWithTransaction
}

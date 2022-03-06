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
	"strconv"
	"time"
)

// startCharging Start charging on the first available Connector. If there is no available Connector, reject the request.
func (cp *ChargePoint) startCharging(tagId string) error {
	if c := cp.connectorManager.FindAvailableConnector(); !util.IsNilInterfaceOrPointer(c) {
		return cp.startChargingConnector(c, tagId)
	}

	return errors.ErrNoAvailableConnectors
}

// startChargingConnector Start charging a connector with the specified ID.
// Send the request to the Central System, turn on and update the status of the Connector,
// start the timer and sample the PowerMeter, if it's enabled.
func (cp *ChargePoint) startChargingConnector(connector connector.Connector, tagId string) error {
	if util.IsNilInterfaceOrPointer(connector) {
		return errors.ErrConnectorNil
	}

	logInfo := cp.logger.WithFields(log.Fields{
		"evseId":      connector.GetEvseId(),
		"connectorId": connector.GetConnectorId(),
		"tagId":       tagId,
	})

	if !connector.IsAvailable() {
		return errors.ErrConnectorUnavailable
	}

	if cp.availability != core.AvailabilityTypeOperative {
		return errors.ErrChargePointUnavailable
	}

	if !cp.isTagAuthorized(tagId) {
		return errors.ErrTagUnauthorized
	}

	request := core.NewStartTransactionRequest(
		connector.GetConnectorId(),
		tagId,
		0,
		types.NewDateTime(time.Now()),
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logInfo.WithError(protoError).Errorf("Server responded with error when starting a transaction")
			return
		}

		startTransactionConf := confirmation.(*core.StartTransactionConfirmation)

		switch startTransactionConf.IdTagInfo.Status {
		case types.AuthorizationStatusAccepted, types.AuthorizationStatusConcurrentTx:
			// Attempt to start charging
			err := connector.StartCharging(strconv.Itoa(startTransactionConf.TransactionId), tagId)
			if err != nil {
				logInfo.WithError(err).Errorf("Unable to start charging connector")
				return
			}

			logInfo.Infof("Started charging connector at %s", time.Now())

			// Schedule timer to stop the transaction at the time limit
			_, err = cp.scheduler.Every(connector.GetMaxChargingTime()).Minutes().LimitRunsTo(1).
				Tag(fmt.Sprintf("connector%dTimer", connector.GetConnectorId())).Do(cp.stopChargingConnector, connector, core.ReasonOther)
			if err != nil {
				logInfo.WithError(err).Errorf("Cannot schedule stop charging")
			}
			break
		case types.AuthorizationStatusBlocked, types.AuthorizationStatusInvalid, types.AuthorizationStatusExpired:
			fallthrough
		default:
			logInfo.Errorf("Transaction unauthorized")
		}
	}

	return util.SendRequest(cp.chargePoint, request, callback)
}

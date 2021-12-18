package v16

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/components/connector"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"log"
	"strconv"
	"time"
)

// startCharging Start charging on the first available ConnectorImpl. If there is no available ConnectorImpl, reject the request.
func (handler *ChargePointHandler) startCharging(tagId string) error {
	if c := handler.connectorManager.FindAvailableConnector(); c != nil {
		return handler.startChargingConnector(c, tagId)
	}

	return chargepoint.ErrNoAvailableConnectors
}

// startChargingConnector Start charging a connector with the specified ID. Send the request to the central system, turn on the ConnectorImpl,
// update the status of the ConnectorImpl, and start the maxChargingTime timer and sample the PowerMeter, if it's enabled.
func (handler *ChargePointHandler) startChargingConnector(connector connector.Connector, tagId string) error {
	if data.IsNilInterfaceOrPointer(connector) {
		return chargepoint.ErrConnectorNil
	}

	if !connector.IsAvailable() {
		return chargepoint.ErrConnectorUnavailable
	}

	if handler.availability != core.AvailabilityTypeOperative {
		return chargepoint.ErrChargePointUnavailable
	}

	if !handler.isTagAuthorized(tagId) {
		return chargepoint.ErrTagUnauthorized
	}

	request := core.NewStartTransactionRequest(
		connector.GetConnectorId(),
		tagId,
		0,
		types2.NewDateTime(time.Now()),
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		startTransactionConf := confirmation.(*core.StartTransactionConfirmation)

		switch startTransactionConf.IdTagInfo.Status {
		case types2.AuthorizationStatusAccepted:
			if startTransactionConf.TransactionId > 0 {
				err := connector.StartCharging(strconv.Itoa(startTransactionConf.TransactionId), tagId)
				if err != nil {
					log.Printf("Unable to start charging connector %d: %s", connector.GetConnectorId(), err)
					return
				}

				log.Printf("Started charging connector %d at %s", connector.GetConnectorId(), time.Now())

				// Schedule timer to stop the transaction at the time limit
				_, err = scheduler.GetScheduler().Every(connector.GetMaxChargingTime()).Minutes().LimitRunsTo(1).
					Tag(fmt.Sprintf("connector%dTimer", connector.GetConnectorId())).Do(handler.stopChargingConnector, connector, core.ReasonOther)
				if err != nil {
					fmt.Println("cannot schedule stop charging:", err)
				}
			}
			break
		case types2.AuthorizationStatusConcurrentTx:
			// todo figure out what to do here exactly
			break
		case types2.AuthorizationStatusBlocked, types2.AuthorizationStatusInvalid, types2.AuthorizationStatusExpired:
			fallthrough
		default:
			log.Printf("Transaction unauthorized at connector %d", connector.GetConnectorId())
		}
	}

	return handler.SendRequest(request, callback)
}

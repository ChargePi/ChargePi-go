package v16

import (
	"fmt"
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
)

func (cp *ChargePoint) StartCharging(evseId, connectorId int, tagId string) error {
	logInfo := cp.logger.WithFields(log.Fields{
		"evseId":      evseId,
		"connectorId": connectorId,
		"tagId":       tagId,
	})
	logInfo.Infof("Starting charging")

	// Charge point must be available to accept transactions
	if cp.availability != core.AvailabilityTypeOperative {
		return chargePoint.ErrChargePointUnavailable
	}

	// Authorize the tag from either the list, cache or the backend.
	if !cp.isTagAuthorized(tagId) {
		return chargePoint.ErrTagUnauthorized
	}

	request := core.NewStartTransactionRequest(
		evseId,
		tagId,
		0,
		types.NewDateTime(time.Now()),
	)

	callback := func(confirmation ocpp.Response, protoError error) {
		if protoError != nil {
			logInfo.WithError(protoError).Warn("Central system responded with an error for %s", confirmation.GetFeatureName())
			return
		}

		startTransactionConf := confirmation.(*core.StartTransactionConfirmation)

		switch startTransactionConf.IdTagInfo.Status {
		case types.AuthorizationStatusAccepted, types.AuthorizationStatusConcurrentTx:
			transactionId := fmt.Sprintf("%d", startTransactionConf.TransactionId)

			err := cp.sessionManager.StartSession(evseId, nil, tagId, transactionId)
			if err != nil {
				logInfo.WithError(err).Error("Unable to start a session")
				return
			}

			// Start the charging on EVSE
			err = cp.evseManager.StartCharging(evseId, nil)
			if err != nil {
				logInfo.WithError(err).Error("Unable to start charging on EVSE")
				return
			}

			logInfo.Infof("Started charging connector at %s", time.Now())
		case types.AuthorizationStatusBlocked, types.AuthorizationStatusInvalid, types.AuthorizationStatusExpired:
			fallthrough
		default:
			logInfo.Warn("Transaction unauthorized")
		}

		// Cache the tag
		err := cp.tagManager.AddTag(tagId, startTransactionConf.IdTagInfo)
		if err != nil {
			logInfo.WithError(err).Warn("Unable to cache tag")
		}
	}

	return cp.sendRequest(request, callback)
}

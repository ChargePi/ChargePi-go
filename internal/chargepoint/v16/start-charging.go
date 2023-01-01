package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"strconv"
	"time"
)

func (cp *ChargePoint) StartCharging(evseId, connectorId int, tagId string) error {
	logInfo := cp.logger.WithFields(log.Fields{
		"evseId":      evseId,
		"connectorId": connectorId,
		"tagId":       tagId,
	})

	if cp.availability != core.AvailabilityTypeOperative {
		return chargePoint.ErrChargePointUnavailable
	}

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
			logInfo.WithError(protoError).Errorf("Central system responded with an error for %s", confirmation.GetFeatureName())
			return
		}

		startTransactionConf := confirmation.(*core.StartTransactionConfirmation)

		switch startTransactionConf.IdTagInfo.Status {
		case types.AuthorizationStatusAccepted, types.AuthorizationStatusConcurrentTx:

			// Attempt to start charging
			err := cp.connectorManager.StartCharging(evseId, tagId, strconv.Itoa(startTransactionConf.TransactionId))
			if err != nil {
				logInfo.WithError(err).Errorf("Unable to start charging connector")
				return
			}

			logInfo.Infof("Started charging connector at %s", time.Now())
		case types.AuthorizationStatusBlocked, types.AuthorizationStatusInvalid, types.AuthorizationStatusExpired:
			fallthrough
		default:
			logInfo.Errorf("Transaction unauthorized")
		}
	}

	return util.SendRequest(cp.chargePoint, request, callback)
}

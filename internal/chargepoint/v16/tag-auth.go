package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"time"
)

// isTagAuthorized Check if the tag is authorized for charging. If the authentication cache is enabled and if it can preauthorize from cache,
// the program will check the cache first and reauthorize with the sendAuthorizeRequest to the central system after 10 seconds.
// If cache is not enabled, it will just execute sendAuthorizeRequest and retrieve the status from the request.
func (cp *ChargePoint) isTagAuthorized(tagId string) bool {
	var (
		response                      = false
		localPreAuthorize, preAuthErr = ocppConfigManager.GetConfigurationValue(v16.LocalPreAuthorize.String())
	)

	if preAuthErr != nil {
		localPreAuthorize = "false"
	}

	if localPreAuthorize == "true" {
		cp.logger.Infof("Authorizing tag %s with cache", tagId)

		tag, err := cp.tagManager.GetTag(tagId)
		if err != nil {
			goto Skip
		}

		switch tag.Status {
		case types.AuthorizationStatusAccepted,
			types.AuthorizationStatusConcurrentTx:
			if tag.ExpiryDate != nil && tag.ExpiryDate.Before(time.Now()) {
				return false
			}

			return true
		case types.AuthorizationStatusInvalid,
			types.AuthorizationStatusBlocked,
			types.AuthorizationStatusExpired:
			return false
		default:
			return false
		}
	}

	// If the card is not in cache or is not authorized, (re)authorize it with the central system
Skip:
	cp.logger.Infof("Authorizing tag %s with central system", tagId)
	tagInfo, err := cp.sendAuthorizeRequest(tagId)
	if err != nil {
		return false
	}

	if tagInfo != nil && tagInfo.Status == types.AuthorizationStatusAccepted {
		response = true
	}

	cp.logger.Debugf("Tag authorization result: %v", response)
	return response
}

// sendAuthorizeRequest Send a AuthorizeRequest to the central system to get information on the tagId status.
// Adds the tag to the cache if it's enabled.
func (cp *ChargePoint) sendAuthorizeRequest(tagId string) (*types.IdTagInfo, error) {
	// Send a request
	response, err := cp.chargePoint.SendRequest(core.NewAuthorizationRequest(tagId))
	if err != nil {
		return nil, err
	}

	authInfo := response.(*core.AuthorizeConfirmation)

	switch authInfo.IdTagInfo.Status {
	case types.AuthorizationStatusBlocked, types.AuthorizationStatusExpired, types.AuthorizationStatusInvalid:
		err = cp.stopChargingConnectorWithTagId(tagId, core.ReasonDeAuthorized)
	}

	_ = cp.tagManager.AddTag(tagId, authInfo.IdTagInfo)

	return authInfo.IdTagInfo, err
}

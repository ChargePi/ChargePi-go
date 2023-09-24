package v16

import (
	"time"

	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/configuration"
)

// isTagAuthorized Check if the tag is authorized for charging. If the authentication cache is enabled and if it can preauthorize from cache,
// the program will check the cache first and reauthorize with the sendAuthorizeRequest to the central system after 10 seconds.
// If cache is not enabled, it will just execute sendAuthorizeRequest and retrieve the status from the request.
func (cp *ChargePoint) isTagAuthorized(tagId string) bool {
	var (
		response             = false
		localPreAuthorize, _ = ocppConfigManager.GetConfigurationValue(v16.LocalPreAuthorize.String())
		logInfo              = cp.logger.WithField("tag", tagId)
	)

	// When local authorization is enabled, get the tag details from the cache or authList before requesting the
	// tag details from the backend.
	if localPreAuthorize != nil && *localPreAuthorize == "true" {
		logInfo.Infof("Authorizing tag %s with cache", tagId)

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
	tagInfo, err := cp.sendAuthorizeRequest(tagId)
	if err != nil {
		logInfo.Warn("Unable to authorize the tag")
		return false
	}

	if tagInfo != nil && tagInfo.Status == types.AuthorizationStatusAccepted {
		response = true
	}

	logInfo.Debugf("Tag authorization result: %v", response)
	return response
}

// sendAuthorizeRequest Send a AuthorizeRequest to the central system to get information on the tag status.
// Adds the tag to the cache and/or localAuthList if it's enabled.
func (cp *ChargePoint) sendAuthorizeRequest(tagId string) (*types.IdTagInfo, error) {
	logInfo := cp.logger.WithField("tag", tagId)
	logInfo.Info("Authorizing the tag with the central system")

	// Authorize the tag with the backend.
	response, err := cp.chargePoint.SendRequest(core.NewAuthorizationRequest(tagId))
	if err != nil {
		logInfo.WithError(err).Error("Tag authorization with the central system failed")

		// An error occurred probably due network issues.
		// If LocalAuthOffline is enabled, try to authenticate from cache or localAuthList.
		localAuthOffline, _ := ocppConfigManager.GetConfigurationValue(v16.LocalAuthorizeOffline.String())
		if localAuthOffline != nil && *localAuthOffline == "true" {
			logInfo.Warn("Offline authorization enabled, getting tag")
			tag, err := cp.tagManager.GetTag(tagId)
			if err != nil {
				return nil, err
			}

			return tag, nil
		}

		return nil, err
	}

	authInfo := response.(*core.AuthorizeConfirmation)
	switch authInfo.IdTagInfo.Status {
	case types.AuthorizationStatusBlocked,
		types.AuthorizationStatusExpired,
		types.AuthorizationStatusInvalid:

		// Stop the transaction and charging process if the StopTransactionOnInvalidId is enabled.
		stopTransactionOnInvalidId, _ := ocppConfigManager.GetConfigurationValue(v16.StopTransactionOnInvalidId.String())
		if stopTransactionOnInvalidId != nil && *stopTransactionOnInvalidId == "true" {
			logInfo.Warn("Tag status invalid or expired, stopping any charging session with the tag")

			sessionId, err := cp.sessionManager.GetSessionWithTagId(tagId)
			if err != nil {
				return nil, err
			}

			err = cp.evseManager.StopCharging(sessionId.EvseId, nil, core.ReasonDeAuthorized)
			if err != nil {
				return nil, err
			}
		}
	}

	// Cache the tag (if available)
	addErr := cp.tagManager.AddTag(tagId, authInfo.IdTagInfo)
	if addErr != nil {
		logInfo.Warn("Unable to add tag to authorization manager")
	}

	return authInfo.IdTagInfo, nil
}

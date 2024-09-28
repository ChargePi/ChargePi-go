package v16

import (
	"time"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// isAuthorized Check if the tag is authorized for charging. If the authentication cache is enabled and
// if it can be preauthorized from cache, check the cache first and reauthorize with the requestTagAuthorization
// to the central system after 10 seconds. If cache is not enabled, the tag info will be fetched from the central system.
func (cp *ChargePoint) isAuthorized(tagId string) bool {
	logInfo := cp.logger.WithField("tag", tagId)

	isAuthorized, err := cp.preAuthorizeFromCache(tagId)
	if err != nil {
		logInfo.WithError(err).Warn("Unable to preauthorize tag from cache")
		return false
	}
	if isAuthorized {
		// Re-request authorization from the central system after 10 seconds.
		return true
	}

	tagInfo, err := cp.requestTagAuthorization(tagId)
	if err != nil {
		logInfo.Warn("Unable to authorize the tag")
		return false
	}

	if tagInfo != nil && tagInfo.Status == types.AuthorizationStatusAccepted {
		return true
	}

	return false
}

// preAuthorizeFromCache Check if the tag is authorized for charging from the cache or localAuthList.
func (cp *ChargePoint) preAuthorizeFromCache(tagId string) (bool, error) {
	logInfo := cp.logger.WithField("tag", tagId)

	localPreAuthorize, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.LocalPreAuthorize)
	if err != nil {
		logInfo.Warn("Unable to determine local authorization configuration")
		return false, err
	}

	// When local authorization is enabled, get the tag details from the cache or authList
	if localPreAuthorize != nil && *localPreAuthorize == "true" {
		logInfo.Infof("Preauthorizing tag %s with cache", tagId)

		tag, err := cp.tagAuthService.GetTag(tagId)
		if err != nil {
			return false, err
		}

		switch tag.Status {
		case types.AuthorizationStatusAccepted,
			types.AuthorizationStatusConcurrentTx:
			if tag.ExpiryDate != nil && tag.ExpiryDate.Before(time.Now()) {
				return false, errors.New("local tag expired")
			}

			return true, nil
		case types.AuthorizationStatusInvalid,
			types.AuthorizationStatusBlocked,
			types.AuthorizationStatusExpired:
			return false, nil
		default:
			return false, nil
		}
	}

	return false, nil
}

// requestTagAuthorization Send a AuthorizeRequest to the central system to get information on the tag status.
// Adds the tag to the cache and/or localAuthList if it's enabled.
func (cp *ChargePoint) requestTagAuthorization(tagId string) (*types.IdTagInfo, error) {
	logInfo := cp.logger.WithField("tag", tagId)
	logInfo.Info("Authorizing the tag with the central system")

	// Authorize the tag with the backend.
	response, err := cp.chargePoint.SendRequest(core.NewAuthorizationRequest(tagId))
	if err != nil {
		logInfo.WithError(err).Error("Tag authorization with the central system failed")
		// An error occurred probably due network issues.
		return cp.authorizeOffline(tagId, logInfo)
	}

	authInfo := response.(*core.AuthorizeConfirmation)
	err2 := cp.checkTagValidity(tagId, authInfo.IdTagInfo.Status)
	if err2 != nil {
		return nil, err2
	}

	// Cache the tag if the cache is enabled.
	addErr := cp.tagAuthService.CacheTag(tagId, authInfo.IdTagInfo)
	if addErr != nil {
		logInfo.Warn("Unable to add tag to authorization manager")
	}

	return authInfo.IdTagInfo, nil
}

// authorizeOffline Authorize the tag from the cache or localAuthList
// if a network error occurs and the offline authorization is enabled.
func (cp *ChargePoint) authorizeOffline(tagId string, logInfo *log.Entry) (*types.IdTagInfo, error) {
	localAuthOffline, _ := cp.settingsManager.GetConfigurationValue(ocpp_v16.LocalAuthorizeOffline)
	if localAuthOffline != nil && *localAuthOffline == "true" {
		logInfo.Warn("Offline authorization enabled, getting tag")
		tag, err := cp.tagAuthService.GetTag(tagId)
		if err != nil {
			return nil, err
		}

		return tag, nil
	}

	return nil, nil
}

// checkTagValidity Check the validity of the tag status and take action if the tag is invalid or expired.
func (cp *ChargePoint) checkTagValidity(tagId string, status types.AuthorizationStatus) error {
	logInfo := cp.logger.WithField("tag", tagId)

	switch status {
	case types.AuthorizationStatusBlocked,
		types.AuthorizationStatusExpired,
		types.AuthorizationStatusInvalid:

		// Stop the transaction and charging process if the StopTransactionOnInvalidId is enabled.
		stopTransactionOnInvalidId, err := cp.settingsManager.GetConfigurationValue(ocpp_v16.StopTransactionOnInvalidId)
		if err != nil {
			return err
		}

		if stopTransactionOnInvalidId != nil && *stopTransactionOnInvalidId == "true" {
			logInfo.Warn("Tag status invalid or expired, stopping any charging session with the tag")

			// todo check the type of error - if not found, proceed without error
			sessionId, err := cp.sessionService.GetSessionWithTagId(tagId)
			if err != nil {
				return err
			}

			err = cp.evseManager.StopCharging(sessionId.EvseId, nil, core.ReasonDeAuthorized)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (cp *ChargePoint) OnClearCache(request *core.ClearCacheRequest) (confirmation *core.ClearCacheConfirmation, err error) {
	cp.logger.Infof("Received request %s", request.GetFeatureName())
	response := core.ClearCacheStatusRejected

	cacheErr := cp.tagAuthService.ClearCache()
	switch {
	case cacheErr == nil:
		cp.logger.Info("Cache cleared")
		response = core.ClearCacheStatusAccepted
	case errors.Is(cacheErr, auth.ErrCacheDisabled):
		cp.logger.Info("Cache not enabled")
		response = core.ClearCacheStatusRejected
	default:
		cp.logger.WithError(cacheErr).Warn("Unable to clear cache")
		response = core.ClearCacheStatusRejected
	}

	return core.NewClearCacheConfirmation(response), nil
}

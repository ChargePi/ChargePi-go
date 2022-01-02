package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	v16 "github.com/xBlaz3kx/ocppManager-go/v16"
	"strconv"
)

// isTagAuthorized Check if the tag is authorized for charging. If the authentication cache is enabled and if it can preauthorize from cache,
// the program will check the cache first and reauthorize with the sendAuthorizeRequest to the central system after 10 seconds.
// If cache is not enabled, it will just execute sendAuthorizeRequest and retrieve the status from the request.
func (cp *ChargePoint) isTagAuthorized(tagId string) bool {
	var (
		response                      = false
		authCacheEnabled, cacheErr    = ocppConfigManager.GetConfigurationValue(v16.AuthorizationCacheEnabled.String())
		localPreAuthorize, preAuthErr = ocppConfigManager.GetConfigurationValue(v16.LocalPreAuthorize.String())
	)

	if cacheErr != nil {
		authCacheEnabled = "false"
	}

	if preAuthErr != nil {
		localPreAuthorize = "false"
	}

	if authCacheEnabled == "true" && localPreAuthorize == "true" {
		log.Infof("Authorizing tag %s with cache", tagId)

		// Check if the tag exists in cache and is valid.
		if cp.authCache.IsTagAuthorized(tagId) {
			// Reauthorize in 10 seconds
			_, schedulerErr := cp.scheduler.Every(10).Seconds().LimitRunsTo(1).Do(cp.sendAuthorizeRequest, tagId)
			if schedulerErr != nil {
				log.WithError(schedulerErr).Errorf("Cannot schedule tag authorization with central system")
			}

			return true
		}
	}

	// If the card is not in cache or is not authorized, (re)authorize it with the central system
	log.Infof("Authorizing tag %s with central system", tagId)
	tagInfo, err := cp.sendAuthorizeRequest(tagId)
	if err != nil {
		return false
	}

	if tagInfo != nil && tagInfo.Status == types.AuthorizationStatusAccepted {
		response = true
	}

	log.Debugf("Tag authorization result: %v", response)
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
		break
	}

	value, err2 := ocppConfigManager.GetConfigurationValue(v16.AuthorizationCacheEnabled.String())
	if err2 == nil && value == "true" {
		cp.authCache.AddTag(tagId, authInfo.IdTagInfo)
	}

	return authInfo.IdTagInfo, err
}

func (cp *ChargePoint) setMaxCachedTags() {
	var (
		maxCachedTagsString, confErr = ocppConfigManager.GetConfigurationValue(v16.LocalAuthListMaxLength.String())
		maxCachedTags, convErr       = strconv.Atoi(maxCachedTagsString)
	)

	if confErr == nil && convErr == nil {
		cp.authCache.SetMaxCachedTags(maxCachedTags)
	}
}

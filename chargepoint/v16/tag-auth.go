package v16

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/components/scheduler"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/conf-manager"
	"github.com/xBlaz3kx/ChargePi-go/data/auth"
	"strconv"
)

// isTagAuthorized Check if the tag is authorized for charging. If the authentication cache is enabled and if it can preauthorize from cache,
// the program will check the cache first and reauthorize with the sendAuthorizeRequest to the central system after 10 seconds.
// If cache is not enabled, it will just execute sendAuthorizeRequest and retrieve the status from the request.
func (handler *ChargePointHandler) isTagAuthorized(tagId string) bool {
	var (
		response                      = false
		authCacheEnabled, cacheErr    = conf_manager.GetConfigurationValue("AuthorizationCacheEnabled")
		localPreAuthorize, preAuthErr = conf_manager.GetConfigurationValue("LocalPreAuthorize")
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
		if auth.IsTagAuthorized(tagId) {
			// Reauthorize in 10 seconds
			_, schedulerErr := scheduler.GetScheduler().Every(10).Seconds().LimitRunsTo(1).Do(handler.sendAuthorizeRequest, tagId)
			if schedulerErr != nil {
				log.Errorf("Cannot schedule check with central system: %v", schedulerErr)
			}

			return true
		}
	}

	// If the card is not in cache or is not authorized, (re)authorize it with the central system
	log.Infof("Authorizing tag %s with central system", tagId)
	tagInfo, err := handler.sendAuthorizeRequest(tagId)
	if err != nil {
		return false
	}

	if tagInfo != nil && tagInfo.Status == types2.AuthorizationStatusAccepted {
		response = true
	}

	log.Infof("Tag authorization result: %v", response)
	return response
}

// sendAuthorizeRequest Send a AuthorizeRequest to the central system to get information on the tagId status.
// Adds the tag to the cache if it's enabled.
func (handler *ChargePointHandler) sendAuthorizeRequest(tagId string) (*types2.IdTagInfo, error) {
	// Send a request
	response, err := handler.chargePoint.SendRequest(core.NewAuthorizationRequest(tagId))
	if err != nil {
		return nil, err
	}

	authInfo := response.(*core.AuthorizeConfirmation)

	switch authInfo.IdTagInfo.Status {
	case types2.AuthorizationStatusBlocked, types2.AuthorizationStatusExpired, types2.AuthorizationStatusInvalid:
		err = handler.stopChargingConnectorWithTagId(tagId, core.ReasonDeAuthorized)
		break
	}

	value, err2 := conf_manager.GetConfigurationValue("AuthorizationCacheEnabled")
	if err2 == nil && value == "true" {
		auth.AddTag(tagId, authInfo.IdTagInfo)
	}

	return authInfo.IdTagInfo, err
}

func (handler *ChargePointHandler) setMaxCachedTags() {
	var (
		maxCachedTagsString, confErr = conf_manager.GetConfigurationValue("MaxCachedTags")
		maxCachedTags, convErr       = strconv.Atoi(maxCachedTagsString)
	)

	if confErr == nil && convErr == nil {
		auth.SetMaxCachedTags(maxCachedTags)
	}
}

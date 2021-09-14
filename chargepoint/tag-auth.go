package chargepoint

import (
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	types2 "github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"strconv"
)

// isTagAuthorized Check if the tag is authorized for charging. If the authentication cache is enabled and if it can preauthorize from cache,
// the program will check the cache first and reauthorize with the sendAuthorizeRequest to the central system after 10 seconds.
// If cache is not enabled, it will just execute sendAuthorizeRequest and retrieve the status from the request.
func (handler *ChargePointHandler) isTagAuthorized(tagId string) bool {
	response := false
	authCacheEnabled, err := settings.GetConfigurationValue("AuthorizationCacheEnabled")
	if err != nil {
		authCacheEnabled = "false"
	}
	localPreAuthorize, err := settings.GetConfigurationValue("LocalPreAuthorize")
	if err != nil {
		localPreAuthorize = "false"
	}
	if authCacheEnabled == "true" && localPreAuthorize == "true" {
		// Check if the tag exists in cache and is valid.
		log.Println("Authorizing tag ", tagId, " with cache")
		if data.IsTagAuthorized(tagId) {
			// Reauthorize in 10 seconds
			_, err2 := scheduler.Every(10).Seconds().LimitRunsTo(1).Do(handler.sendAuthorizeRequest, tagId)
			if err2 != nil {
				log.Println(err2)
			}
			return true
		}
	}
	//If the card is not in cache or is not authorized, (re)authorize it with the central system
	log.Println("Authorizing tag with central system: ", tagId)
	tagInfo, err := handler.sendAuthorizeRequest(tagId)
	if tagInfo != nil && tagInfo.Status == types2.AuthorizationStatusAccepted {
		response = true
	}
	log.Println("Tag authorization result: ", response)
	return response
}

// sendAuthorizeRequest Send a AuthorizeRequest to the central system to get information on the tagId status.
// Adds the tag to the cache if it's enabled.
func (handler *ChargePointHandler) sendAuthorizeRequest(tagId string) (*types2.IdTagInfo, error) {
	var err error
	handler.mu.Lock()
	response, err := handler.chargePoint.SendRequest(core.AuthorizeRequest{IdTag: tagId})
	handler.mu.Unlock()
	authInfo := response.(*core.AuthorizeConfirmation)
	switch authInfo.IdTagInfo.Status {
	case types2.AuthorizationStatusBlocked, types2.AuthorizationStatusExpired, types2.AuthorizationStatusInvalid:
		err = handler.stopChargingConnectorWithTagId(tagId, core.ReasonDeAuthorized)
		break
	}
	value, err2 := settings.GetConfigurationValue("AuthorizationCacheEnabled")
	if err2 == nil && value == "true" {
		data.AddTag(tagId, authInfo.IdTagInfo)
	}
	return authInfo.IdTagInfo, err
}

func (handler *ChargePointHandler) setMaxCachedTags() {
	maxCachedTagsString, err := settings.GetConfigurationValue("MaxCachedTags")
	maxCachedTags, err := strconv.Atoi(maxCachedTagsString)
	if err == nil {
		data.SetMaxCachedTags(maxCachedTags)
	}
}

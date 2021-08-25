package data

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"sync"
	"time"
)

type AuthorizationModule struct {
	Version       int
	MaxCachedTags int
	tags          []types.IdTagInfo
}

var AuthCache *goCache.Cache

//init Read the authorization persistence file.
func init() {
	var auth *AuthorizationModule
	once := sync.Once{}
	once.Do(func() {
		AuthCache = goCache.New(time.Minute*10, time.Minute*10)
		settings.DecodeFile("configs/auth.json", &auth)
		AuthCache.Set("AuthCacheVersion", auth.Version, goCache.NoExpiration)
		AuthCache.Set("AuthCacheMaxTags", 0, goCache.NoExpiration)
		for _, tag := range auth.tags {
			if tag.ExpiryDate != nil {
				AuthCache.Set(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag, tag.ExpiryDate.Sub(time.Now()))
				continue
			}
			AuthCache.SetDefault(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag)
		}
		log.Printf("Read auth file version %d with tags %s", auth.Version, auth.tags)
	})
}

// AddTag Add a tag to the global authorization cache.
func AddTag(tagId string, tagInfo *types.IdTagInfo) {
	cacheMaxTags, isFound := AuthCache.Get("AuthCacheMaxTags")
	if !isFound {
		return
	}
	var maxTags = cacheMaxTags.(int)
	if AuthCache.ItemCount()-2 >= maxTags {
		return
	}
	if tagInfo.ExpiryDate != nil {
		AuthCache.Set(fmt.Sprintf("AuthTag%s", tagId), tagInfo, tagInfo.ExpiryDate.Sub(time.Now()))
		return
	}
	AuthCache.SetDefault(fmt.Sprintf("AuthTag%s", tagId), tagInfo)
}

// RemoveTag Remove a tag from the global authorization cache.
func RemoveTag(tagId string) {
	AuthCache.Delete(fmt.Sprintf("AuthTag%s", tagId))
}

// RemoveCachedTags Remove all tags from the global authorization cache.
func RemoveCachedTags() {
	version, isVersionFound := AuthCache.Get("AuthCacheVersion")
	maxCachedTags, isMaxFound := AuthCache.Get("AuthCacheMaxTags")
	AuthCache.Flush()
	if isVersionFound {
		AuthCache.Set("AuthCacheVersion", version, goCache.NoExpiration)
	}
	if isMaxFound {
		AuthCache.Set("AuthCacheMaxTags", maxCachedTags, goCache.NoExpiration)
	}
}

// SetMaxCachedTags Set the maximum number of tags allowed in the global authorization cache.
func SetMaxCachedTags(number int) {
	if number > 0 {
		AuthCache.Set("AuthCacheMaxTags", number, goCache.NoExpiration)
	}
}

// IsTagAuthorized Check if the tag exists in the global authorization cache, the status of the tag is "Accepted" and if it has not expired yet.
func IsTagAuthorized(tagId string) bool {
	log.Println("Checking if tag authorized", tagId)
	tagObject, isFound := AuthCache.Get(fmt.Sprintf("AuthTag%s", tagId))
	if isFound {
		tagInfo := tagObject.(*types.IdTagInfo)
		if tagInfo.Status == types.AuthorizationStatusAccepted {
			if tagInfo.ExpiryDate != nil && tagInfo.ExpiryDate.Before(time.Now()) {
				return false
			}
			log.Println("Tag ", tagId, " authorized with cache")
			return true
		}
		return false
	}
	return false
}

package auth

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/components/cache"
	"github.com/xBlaz3kx/ChargePi-go/components/settings"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type (
	AuthorizationFile struct {
		Version       int               `fig:"Version" validation:"required"`
		MaxCachedTags int               `fig:"MaxCachedTags" validation:"required"`
		Tags          []types.IdTagInfo `fig:"Tags"`
	}
)

var authCache *goCache.Cache

//init Read the authorization file.
func init() {
	once := sync.Once{}
	once.Do(func() {
		GetAuthCache()
	})
}

func GetAuthCache() *goCache.Cache {
	if authCache == nil {
		log.Info("Creating auth cache..")
		authCache = goCache.New(time.Minute*10, time.Minute*10)
	}

	return authCache
}

// LoadAuthFile loads tags from the cache file
func LoadAuthFile() {
	var (
		auth         AuthorizationFile
		authFilePath = ""
		err          error
	)

	authPath, isFound := cache.GetCache().Get("authFilePath")
	if isFound {
		authFilePath = authPath.(string)
	}

	err = fig.Load(&auth,
		fig.File(filepath.Base(authFilePath)),
		fig.Dirs(filepath.Dir(authFilePath)))
	if err != nil {
		//todo temporary fix - tags with ExpiryDate won't unmarshall successfully
		log.Errorf("Unable to load auth file: %v", err)
	}

	authCache.Set("AuthCacheVersion", auth.Version, goCache.NoExpiration)
	authCache.Set("AuthCacheMaxTags", auth.MaxCachedTags, goCache.NoExpiration)

	if auth.Tags != nil {
		for _, tag := range auth.Tags {
			log.Debugf("Adding tag: %v", tag)
			if tag.ExpiryDate != nil {
				authCache.Set(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag, tag.ExpiryDate.Sub(time.Now()))
				continue
			}

			authCache.SetDefault(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag)
		}
	}

	log.Infof("Read auth file version %d with Tags %s", auth.Version, auth.Tags)
}

// AddTag Add a tag to the global authorization cache.
func AddTag(tagId string, tagInfo *types.IdTagInfo) {
	var (
		maxTags        int
		expirationTime = time.Minute * 10
	)

	cacheMaxTags, isFound := authCache.Get("AuthCacheMaxTags")
	if !isFound {
		return
	}
	maxTags = cacheMaxTags.(int)

	if authCache.ItemCount()-2 >= maxTags {
		return
	}

	if tagInfo.ExpiryDate != nil {
		expirationTime = tagInfo.ExpiryDate.Sub(time.Now())
	}

	// Add a tag if it doesn't exist in the cache already
	err := authCache.Add(fmt.Sprintf("AuthTag%s", tagId), *tagInfo, expirationTime)
	if err != nil {
		log.Errorf("Error adding tag to cache: %v", err)
	}
}

// RemoveTag Remove a tag from the global authorization cache.
func RemoveTag(tagId string) {
	authCache.Delete(fmt.Sprintf("AuthTag%s", tagId))
}

// RemoveCachedTags Remove all Tags from the global authorization cache.
func RemoveCachedTags() {
	log.Info("Flushing auth cache")
	var (
		version, isVersionFound   = authCache.Get("AuthCacheVersion")
		maxCachedTags, isMaxFound = authCache.Get("AuthCacheMaxTags")
	)

	authCache.Flush()

	if !isVersionFound {
		version = 1
	}

	if !isMaxFound {
		maxCachedTags = 0
	}

	authCache.Set("AuthCacheVersion", version, goCache.NoExpiration)
	authCache.Set("AuthCacheMaxTags", maxCachedTags, goCache.NoExpiration)
}

// SetMaxCachedTags Set the maximum number of Tags allowed in the global authorization cache.
func SetMaxCachedTags(number int) {
	if number > 0 {
		log.Infof("Set max cached tags to %d", number)
		authCache.Set("AuthCacheMaxTags", number, goCache.NoExpiration)
	}
}

func DumpTags() {
	log.Infof("Writing tags to file..")
	var (
		authTags                  []types.IdTagInfo
		authFilePath, isFound     = cache.Cache.Get("authFilePath")
		version, isVersionFound   = authCache.Get("AuthCacheVersion")
		maxCachedTags, isMaxFound = authCache.Get("AuthCacheMaxTags")
	)

	if !isFound {
		return
	}

	if !isVersionFound {
		version = 1
	}

	if !isMaxFound {
		maxCachedTags = 0
	}

	for key, item := range authCache.Items() {
		if strings.Contains(key, "AuthTag") && !item.Expired() {
			authTags = append(authTags, item.Object.(types.IdTagInfo))
		}
	}

	err := settings.WriteToFile(authFilePath.(string), AuthorizationFile{
		Version:       version.(int),
		MaxCachedTags: maxCachedTags.(int),
		Tags:          authTags,
	})
	if err != nil {
		log.Errorf("Error updating auth cache file: %v", err)
	}
}

// IsTagAuthorized Check if the tag exists in the global authorization cache, the status of the tag is "Accepted" and if it has not expired yet.
func IsTagAuthorized(tagId string) bool {
	log.Infof("Checking if tag authorized %s", tagId)

	tagObject, isFound := authCache.Get(fmt.Sprintf("AuthTag%s", tagId))
	if isFound {
		tagInfo := tagObject.(types.IdTagInfo)
		switch tagInfo.Status {
		case types.AuthorizationStatusAccepted,
			types.AuthorizationStatusConcurrentTx:
			if tagInfo.ExpiryDate != nil && tagInfo.ExpiryDate.Before(time.Now()) {
				return false
			}

			log.Infof("Tag %s authorized with cache", tagId)
			return true
		case types.AuthorizationStatusInvalid,
			types.AuthorizationStatusBlocked,
			types.AuthorizationStatusExpired:
			return false
		default:
			return false
		}
	}

	return false
}

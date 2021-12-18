package auth

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/components/cache"
	"github.com/xBlaz3kx/ChargePi-go/components/settings"
	"log"
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
		log.Println(err)
	}

	authCache.Set("AuthCacheVersion", auth.Version, goCache.NoExpiration)
	authCache.Set("AuthCacheMaxTags", auth.MaxCachedTags, goCache.NoExpiration)

	if auth.Tags != nil {
		for _, tag := range auth.Tags {
			log.Println("Adding tag:", tag)
			if tag.ExpiryDate != nil {
				authCache.Set(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag, tag.ExpiryDate.Sub(time.Now()))
				continue
			}

			authCache.SetDefault(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag)
		}
	}

	log.Printf("Read auth file version %d with Tags %s", auth.Version, auth.Tags)
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
		log.Println(err)
	}
}

// RemoveTag Remove a tag from the global authorization cache.
func RemoveTag(tagId string) {
	authCache.Delete(fmt.Sprintf("AuthTag%s", tagId))
}

// RemoveCachedTags Remove all Tags from the global authorization cache.
func RemoveCachedTags() {
	version, isVersionFound := authCache.Get("AuthCacheVersion")
	maxCachedTags, isMaxFound := authCache.Get("AuthCacheMaxTags")
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
		authCache.Set("AuthCacheMaxTags", number, goCache.NoExpiration)
	}
}

func DumpTags() {
	log.Println("Writing tags to file..")
	authFilePath, isFound := cache.Cache.Get("authFilePath")
	if !isFound {
		return
	}

	var authTags []types.IdTagInfo
	for key, item := range authCache.Items() {
		if strings.Contains(key, "AuthTag") && !item.Expired() {
			authTags = append(authTags, item.Object.(types.IdTagInfo))
		}
	}

	version, isVersionFound := authCache.Get("AuthCacheVersion")
	maxCachedTags, isMaxFound := authCache.Get("AuthCacheMaxTags")

	if !isVersionFound {
		version = 1
	}

	if !isMaxFound {
		maxCachedTags = 0
	}

	err := settings.WriteToFile(authFilePath.(string), AuthorizationFile{
		Version:       version.(int),
		MaxCachedTags: maxCachedTags.(int),
		Tags:          authTags,
	})
	if err != nil {
		log.Println(err)
	}
}

// IsTagAuthorized Check if the tag exists in the global authorization cache, the status of the tag is "Accepted" and if it has not expired yet.
func IsTagAuthorized(tagId string) bool {
	log.Println("Checking if tag authorized", tagId)
	tagObject, isFound := authCache.Get(fmt.Sprintf("AuthTag%s", tagId))
	if isFound {
		tagInfo := tagObject.(types.IdTagInfo)
		switch tagInfo.Status {
		case types.AuthorizationStatusAccepted,
			types.AuthorizationStatusConcurrentTx:
			if tagInfo.ExpiryDate != nil && tagInfo.ExpiryDate.Before(time.Now()) {
				return false
			}

			log.Println("Tag ", tagId, " authorized with cache")
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

package data

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AuthorizationModule struct {
	Version       int               `fig:"Version" validation:"required"`
	MaxCachedTags int               `fig:"MaxCachedTags" validation:"required"`
	Tags          []types.IdTagInfo `fig:"Tags"`
}

var AuthCache *goCache.Cache

//init Read the authorization persistence file.
func init() {
	once := sync.Once{}
	once.Do(func() {
		AuthCache = goCache.New(time.Minute*10, time.Minute*10)
	})
}

// GetAuthFile read Tags from the persistence cache.
func GetAuthFile() {
	var (
		auth         AuthorizationModule
		authFilePath = ""
		err          error
	)
	authPath, isFound := cache.Cache.Get("authFilePath")
	if isFound {
		authFilePath = authPath.(string)
	}

	err = fig.Load(&auth,
		fig.File(filepath.Base(authFilePath)),
		fig.Dirs(filepath.Dir(authFilePath)))
	if err != nil {
		//log.Fatal(err)
		//todo temporary fix - tags with ExpiryDate won't unmarshall successfully
		log.Println(err)
	}

	AuthCache.Set("AuthCacheVersion", auth.Version, goCache.NoExpiration)
	AuthCache.Set("AuthCacheMaxTags", auth.MaxCachedTags, goCache.NoExpiration)

	if auth.Tags != nil {

		for _, tag := range auth.Tags {
			log.Println("Adding tag:", tag)
			if tag.ExpiryDate != nil {
				AuthCache.Set(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag, tag.ExpiryDate.Sub(time.Now()))
				continue
			}
			AuthCache.SetDefault(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag)
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
	cacheMaxTags, isFound := AuthCache.Get("AuthCacheMaxTags")
	if !isFound {
		return
	}
	maxTags = cacheMaxTags.(int)

	if AuthCache.ItemCount()-2 >= maxTags {
		return
	}

	if tagInfo.ExpiryDate != nil {
		expirationTime = tagInfo.ExpiryDate.Sub(time.Now())
	}
	// add a tag if it doesn't exist in the cache already
	err := AuthCache.Add(fmt.Sprintf("AuthTag%s", tagId), *tagInfo, expirationTime)
	if err != nil {
		log.Println(err)
	}
}

// RemoveTag Remove a tag from the global authorization cache.
func RemoveTag(tagId string) {
	AuthCache.Delete(fmt.Sprintf("AuthTag%s", tagId))
}

// RemoveCachedTags Remove all Tags from the global authorization cache.
func RemoveCachedTags() {
	version, isVersionFound := AuthCache.Get("AuthCacheVersion")
	maxCachedTags, isMaxFound := AuthCache.Get("AuthCacheMaxTags")
	AuthCache.Flush()

	if !isVersionFound {
		version = 1
	}
	AuthCache.Set("AuthCacheVersion", version, goCache.NoExpiration)

	if !isMaxFound {
		maxCachedTags = 0
	}
	AuthCache.Set("AuthCacheMaxTags", maxCachedTags, goCache.NoExpiration)
}

// SetMaxCachedTags Set the maximum number of Tags allowed in the global authorization cache.
func SetMaxCachedTags(number int) {
	if number > 0 {
		AuthCache.Set("AuthCacheMaxTags", number, goCache.NoExpiration)
	}
}

func DumpTags() {
	log.Println("Writing tags to file..")
	authFilePath, isFound := cache.Cache.Get("authFilePath")
	if !isFound {
		return
	}

	var authTags []types.IdTagInfo
	for key, item := range AuthCache.Items() {
		if strings.Contains(key, "AuthTag") && !item.Expired() {
			authTags = append(authTags, item.Object.(types.IdTagInfo))
		}
	}

	version, isVersionFound := AuthCache.Get("AuthCacheVersion")
	maxCachedTags, isMaxFound := AuthCache.Get("AuthCacheMaxTags")

	if !isVersionFound {
		version = 1
	}

	if !isMaxFound {
		maxCachedTags = 0
	}

	err := settings.WriteToFile(authFilePath.(string), AuthorizationModule{
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

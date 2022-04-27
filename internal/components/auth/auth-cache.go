package auth

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/components/settings"
	settingsData "github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"path/filepath"
	"strings"
	"time"
)

const (
	VersionKey = "AuthCacheVersion"
	MaxTagsKey = "AuthCacheMaxTags"
)

type (
	Cache struct {
		cache    *goCache.Cache
		filePath string
	}
)

func NewAuthCache(filePath string) *Cache {
	cache := goCache.New(time.Minute*10, time.Minute*10)
	// Defaults
	cache.Set(VersionKey, 1, goCache.NoExpiration)
	cache.Set(MaxTagsKey, 0, goCache.NoExpiration)

	return &Cache{
		cache:    cache,
		filePath: filePath,
	}
}

// LoadAuthFile loads tags from the cache file
func (c *Cache) LoadAuthFile() {
	var (
		auth settingsData.AuthorizationFile
		err  error
	)

	err = fig.Load(&auth,
		fig.File(filepath.Base(c.filePath)),
		fig.Dirs(filepath.Dir(c.filePath)))
	if err != nil {
		//todo temporary fix - tags with ExpiryDate won't unmarshall successfully
		log.WithError(err).Errorf("Unable to load authorization file")
	}

	c.cache.Set(VersionKey, auth.Version, goCache.NoExpiration)
	c.cache.Set(MaxTagsKey, auth.MaxCachedTags, goCache.NoExpiration)
	loadTags(c.cache, auth.Tags)

	log.Infof("Read auth file version %d with Tags %s", auth.Version, auth.Tags)
}

// AddTag Add a tag to the global authorization cache.
func (c *Cache) AddTag(tagId string, tagInfo *types.IdTagInfo) {
	var (
		maxTags        int
		expirationTime = time.Minute * 10
	)

	cacheMaxTags, isFound := c.cache.Get(MaxTagsKey)
	if !isFound {
		maxTags = 0
	}
	maxTags = cacheMaxTags.(int)

	if c.cache.ItemCount() >= maxTags+2 {
		return
	}

	if tagInfo.ExpiryDate != nil {
		expirationTime = tagInfo.ExpiryDate.Sub(time.Now())
	}

	// Add a tag if it doesn't exist in the cache already
	err := c.cache.Add(fmt.Sprintf("AuthTag%s", tagId), *tagInfo, expirationTime)
	if err != nil {
		log.WithError(err).Errorf("Error adding tag to cache")
	}
}

// RemoveTag Remove a tag from the global authorization cache.
func (c *Cache) RemoveTag(tagId string) {
	c.cache.Delete(fmt.Sprintf("AuthTag%s", tagId))
}

// RemoveCachedTags Remove all Tags from the global authorization cache.
func (c *Cache) RemoveCachedTags() {
	log.Debugf("Flushing auth cache")
	var (
		version, isVersionFound   = c.cache.Get(VersionKey)
		maxCachedTags, isMaxFound = c.cache.Get(MaxTagsKey)
	)

	c.cache.Flush()

	if !isVersionFound {
		version = 1
	}

	if !isMaxFound {
		maxCachedTags = 0
	}

	// Reset the version and max tags
	c.cache.Set(VersionKey, version, goCache.NoExpiration)
	c.cache.Set(MaxTagsKey, maxCachedTags, goCache.NoExpiration)
}

// SetMaxCachedTags Set the maximum number of Tags allowed in the global authorization cache.
func (c *Cache) SetMaxCachedTags(number int) {
	if number > 0 {
		log.Debugf("Set max cached tags to %d", number)
		c.cache.Set(MaxTagsKey, number, goCache.NoExpiration)
	}
}

func (c *Cache) DumpTags() {
	log.Debug("Writing tags to file..")
	var (
		authTags                  []types.IdTagInfo
		version, isVersionFound   = c.cache.Get(VersionKey)
		maxCachedTags, isMaxFound = c.cache.Get(MaxTagsKey)
	)

	if !isVersionFound {
		version = 1
	}

	if !isMaxFound {
		maxCachedTags = 0
	}

	for key, item := range c.cache.Items() {
		if strings.Contains(key, "AuthTag") && !item.Expired() {
			authTags = append(authTags, item.Object.(types.IdTagInfo))
		}
	}

	err := settings.WriteToFile(c.filePath, settingsData.AuthorizationFile{
		Version:       version.(int),
		MaxCachedTags: maxCachedTags.(int),
		Tags:          authTags,
	})
	if err != nil {
		log.WithError(err).Errorf("Error updating auth cache file")
	}
}

// IsTagAuthorized Check if the tag exists in the global authorization cache, the status of the tag is "Accepted" and if it has not expired yet.
func (c *Cache) IsTagAuthorized(tagId string) bool {
	log.Infof("Checking if tag authorized %s", tagId)

	tagObject, isFound := c.cache.Get(fmt.Sprintf("AuthTag%s", tagId))
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

// loadTags loads the tags into the cache
func loadTags(cache *goCache.Cache, tags []types.IdTagInfo) {
	if tags != nil {
		for _, tag := range tags {
			log.Tracef("Adding tag: %v", tag)
			if tag.ExpiryDate != nil {
				cache.Set(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag, tag.ExpiryDate.Sub(time.Now()))
				continue
			}

			cache.SetDefault(fmt.Sprintf("AuthTag%s", tag.ParentIdTag), tag)
		}
	}
}

package auth

import (
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	Cache interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo)
		GetTag(tagId string) (*types.IdTagInfo, error)
		SetMaxCachedTags(number int)
		RemoveCachedTags()
		LoadFromFile() error
		WriteToFile() error
	}

	CacheImpl struct {
		cache   *goCache.Cache
		maxTags int
	}
)

func NewAuthCache() *CacheImpl {
	return &CacheImpl{
		cache:   goCache.New(goCache.NoExpiration, time.Minute*3),
		maxTags: 0,
	}
}

// AddTag Add a tag to the authorization cache.
func (c *CacheImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) {
	logInfo := log.WithFields(log.Fields{
		"tagId":      tagId,
		"tagDetails": tagInfo,
	})
	expirationTime := goCache.NoExpiration

	// Check if the cache is not full
	if c.cache.ItemCount() >= c.maxTags {
		return
	}

	if tagInfo.ExpiryDate != nil {
		expirationTime = tagInfo.ExpiryDate.Sub(time.Now())
	}

	// Add a tag if it doesn't exist in the cache already
	logInfo.Debug("Adding a tag to cache")
	err := c.cache.Add(fmt.Sprintf("AuthTag%s", tagId), *tagInfo, expirationTime)
	if err != nil {
		logInfo.WithError(err).Errorf("Error adding tag to cache")
		return
	}

	defer func(c *CacheImpl) {
		err := c.WriteToFile()
		if err != nil {

		}
	}(c)
}

// RemoveTag Remove a tag from the authorization cache.
func (c *CacheImpl) RemoveTag(tagId string) {
	log.WithField("tagId", tagId).Debug("Removing a tag from cache")
	c.cache.Delete(fmt.Sprintf("AuthTag%s", tagId))
	defer func(c *CacheImpl) {
		err := c.WriteToFile()
		if err != nil {

		}
	}(c)
}

// RemoveCachedTags Remove all Tags from the authorization cache.
func (c *CacheImpl) RemoveCachedTags() {
	log.Debugf("Flushing auth cache")
	c.cache.Flush()

	err := c.WriteToFile()
	if err != nil {

	}
}

// SetMaxCachedTags Set the maximum number of Tags allowed in the authorization cache.
func (c *CacheImpl) SetMaxCachedTags(number int) {
	if number > 0 {
		log.Debugf("Set max cached tags to %d", number)
		c.maxTags = number
	}
}

// GetTag Get a tag from cache
func (c *CacheImpl) GetTag(tagId string) (*types.IdTagInfo, error) {
	log.WithField("tagId", tagId).Info("Getting a tag from cache")

	tagObject, isFound := c.cache.Get(fmt.Sprintf("AuthTag%s", tagId))
	if isFound {
		tagInfo := tagObject.(types.IdTagInfo)
		return &tagInfo, nil
	}

	return nil, ErrTagNotFound
}

func (c *CacheImpl) LoadFromFile() error {
	log.Info("Loading auth cache from file")
	return nil
}

func (c *CacheImpl) WriteToFile() error {
	log.Debug("Writing auth cache to a file")
	return nil
}

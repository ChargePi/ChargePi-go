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
	}

	CacheImpl struct {
		cache   *goCache.Cache
		maxTags int
	}
)

func NewAuthCache() *CacheImpl {
	return &CacheImpl{
		cache:   goCache.New(time.Minute*10, time.Minute*10),
		maxTags: 0,
	}
}

// AddTag Add a tag to the authorization cache.
func (c *CacheImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) {
	var (
		expirationTime = time.Minute * 10
	)

	if c.cache.ItemCount() >= c.maxTags+2 {
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

// RemoveTag Remove a tag from the authorization cache.
func (c *CacheImpl) RemoveTag(tagId string) {
	c.cache.Delete(fmt.Sprintf("AuthTag%s", tagId))
}

// RemoveCachedTags Remove all Tags from the authorization cache.
func (c *CacheImpl) RemoveCachedTags() {
	log.Debugf("Flushing auth cache")
	c.cache.Flush()
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
	log.Infof("Fetching the tag", tagId)

	tagObject, isFound := c.cache.Get(fmt.Sprintf("AuthTag%s", tagId))
	if isFound {
		tagInfo := tagObject.(types.IdTagInfo)
		return &tagInfo, nil
	}

	return nil, ErrTagNotFound
}

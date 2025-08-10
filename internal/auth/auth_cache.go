package auth

import (
	"errors"

	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/lorenzodonini/ocpp-go/ocppj"
	log "github.com/sirupsen/logrus"
)

type (
	Cache interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo) error
		GetTag(tagId string) (*types.IdTagInfo, error)
		SetMaxCachedTags(number int)
		RemoveCachedTags() error
		RemoveTag(tagId string) error
	}

	cache struct {
		repository TagRepository
		maxTags    int
		logger     log.FieldLogger
	}
)

type TagRepository interface {
	AddTag(tagId string, tagInfo *types.IdTagInfo) error
	RemoveTag(tagId string) error
	GetTag(tagId string) (*types.IdTagInfo, error)
	GetTags() ([]*types.IdTagInfo, error)
	RemoveAllTags() error
	// RemoveExpiredTags() error
}

func newAuthCache(repository TagRepository) *cache {
	return &cache{
		repository: repository,
		maxTags:    0,
		logger:     log.StandardLogger().WithField("component", "auth-cache"),
		// todo scheduler to periodically clean up the cache?
	}
}

// AddTag Add a tag to the authorization cache.
func (c *cache) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	logInfo := c.logger.WithField("tagId", tagId)
	logInfo.Debug("Adding a tag to cache")

	if stringUtils.IsEmpty(tagId) {
		return ErrInvalidTagId
	}

	if util.IsNilInterfaceOrPointer(tagInfo) {
		return ErrTagNil
	}

	// Validate
	err := ocppj.Validate.Struct(tagInfo)
	if err != nil {
		return err
	}

	tags, err := c.repository.GetTags()
	if err != nil {
		logInfo.WithError(err).Error("Error getting tags from cache")
		return err
	}

	if len(tags) >= c.maxTags {
		return errors.New("cache is full")
	}

	// Add a expectedTag if it doesn't exist in the cache.
	return c.repository.AddTag(tagId, tagInfo)
}

// RemoveTag Remove a tag from the authorization cache.
func (c *cache) RemoveTag(tagId string) error {
	logInfo := c.logger.WithField("tagId", tagId)
	logInfo.Debug("Removing a tag from cache")

	// Remove a expectedTag if it exists in the cache.
	err := c.repository.RemoveTag(tagId)
	if err != nil {
		logInfo.WithError(err).Error("Error removing tag from cache")
		return err
	}

	return nil
}

// RemoveCachedTags Remove all tags from the authorization cache.
func (c *cache) RemoveCachedTags() error {
	c.logger.Debugf("Flushing auth cache")

	// Remove all cached keys from database
	err := c.repository.RemoveAllTags()
	if err != nil {
		c.logger.WithError(err).Error("Error flushing auth cache")
		return err
	}

	return nil
}

// SetMaxCachedTags Set the maximum number of tags allowed in the authorization cache.
func (c *cache) SetMaxCachedTags(number int) {
	c.logger.Debugf("Set max cached tags to %d", number)

	if number > 0 {
		c.maxTags = number
	}
}

// GetTag Get a tag with id from the authorization cache.
func (c *cache) GetTag(tagId string) (*types.IdTagInfo, error) {
	logInfo := c.logger.WithField("tagId", tagId)
	logInfo.Info("Getting a tag from cache")

	return c.repository.GetTag(tagId)
}

package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

type (
	Cache interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo)
		GetTag(tagId string) (*types.IdTagInfo, error)
		SetMaxCachedTags(number int)
		RemoveCachedTags()
	}

	BadgerCache struct {
		db      *badger.DB
		maxTags int
		logger  *zap.Logger
	}
)

func NewAuthCache(logger *zap.Logger, db *badger.DB) *BadgerCache {
	return &BadgerCache{
		db:      db,
		maxTags: 0,
		logger:  logger.Named("auth_cache"),
	}
}

func getTagKey(tagId string) []byte {
	return []byte(fmt.Sprintf("cached-tag-%s", tagId))
}

// AddTag Add a tag to the authorization cache.
func (c *BadgerCache) AddTag(tagId string, tagInfo *types.IdTagInfo) {
	logInfo := c.logger.With(zap.String("tagId", tagId))
	logInfo.Debug("Adding a tag to cache")

	// Add a tag if it doesn't exist in the cache.
	err := c.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(getTagKey(tagId))
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}

		authTag := getTag(tagId, tagInfo)
		if authTag != nil {
			return nil
		}

		err = txn.Set(getTagKey(tagId), authTag)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
	if err != nil {
		logInfo.With(zap.Error(err)).Error("Error adding tag to cache")
		return
	}
}

// RemoveTag Remove a tag from the authorization cache.
func (c *BadgerCache) RemoveTag(tagId string) {
	logger := c.logger.With(zap.String("tagId", tagId))
	logger.Debug("Removing a tag from cache")

	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getTagKey(tagId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("Error removing tag from cache")
	}
}

// RemoveCachedTags Remove all Tags from the authorization cache.
func (c *BadgerCache) RemoveCachedTags() {
	c.logger.Debug("Flushing auth cache")

	// Remove all cached keys from database
	err := c.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := getTagKey("")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := txn.Delete(item.Key())
			if err != nil {
				return err
			}
		}
		return txn.Commit()
	})
	if err != nil {
		c.logger.With(zap.Error(err)).Error("Error flushing auth cache")
	}
}

// SetMaxCachedTags Set the maximum number of Tags allowed in the authorization cache.
func (c *BadgerCache) SetMaxCachedTags(number int) {
	c.logger.Sugar().Debugf("Set max cached tags to %d", number)

	if number > 0 {
		c.maxTags = number
	}
}

// GetTag Get a tag from cache based on the tag ID.
func (c *BadgerCache) GetTag(tagId string) (*types.IdTagInfo, error) {
	logger := c.logger.With(zap.String("tagId", tagId))
	logger.Info("Getting a tag from cache")

	var tagInfo localauth.AuthorizationData

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getTagKey(tagId))
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}

		var tagCopy []byte
		_, err = item.ValueCopy(tagCopy)
		if err != nil {
			return err
		}

		err = json.Unmarshal(tagCopy, &tagInfo)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("Error getting tag from cache")
		return nil, err
	}

	return tagInfo.IdTagInfo, nil
}

// getTag transforms a tag struct into a byte array.
func getTag(tagId string, tagInfo *types.IdTagInfo) []byte {
	authTag := localauth.AuthorizationData{
		IdTag:     tagId,
		IdTagInfo: tagInfo,
	}

	tag, err := json.Marshal(authTag)
	if err != nil {
		return nil
	}

	return tag
}

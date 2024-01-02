package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
)

type (
	Cache interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo)
		GetTag(tagId string) (*types.IdTagInfo, error)
		SetMaxCachedTags(number int)
		RemoveCachedTags()
	}

	CacheImpl struct {
		db      *badger.DB
		maxTags int
		logger  log.FieldLogger
	}
)

func NewAuthCache(db *badger.DB) *CacheImpl {
	return &CacheImpl{
		db:      db,
		maxTags: 0,
		logger:  log.StandardLogger().WithField("component", "auth-cache"),
	}
}

func getTagKey(tagId string) []byte {
	return []byte(fmt.Sprintf("cached-tag-%s", tagId))
}

// AddTag Add a tag to the authorization cache.
func (c *CacheImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) {
	logInfo := c.logger.WithField("tagId", tagId)
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
		logInfo.WithError(err).Error("Error adding tag to cache")
		return
	}
}

// RemoveTag Remove a tag from the authorization cache.
func (c *CacheImpl) RemoveTag(tagId string) {
	logInfo := c.logger.WithField("tagId", tagId)
	logInfo.Debug("Removing a tag from cache")

	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getTagKey(tagId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
	if err != nil {
		logInfo.WithError(err).Error("Error removing tag from cache")
	}
}

// RemoveCachedTags Remove all Tags from the authorization cache.
func (c *CacheImpl) RemoveCachedTags() {
	log.Debugf("Flushing auth cache")

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
		log.WithError(err).Error("Error flushing auth cache")
	}
}

// SetMaxCachedTags Set the maximum number of Tags allowed in the authorization cache.
func (c *CacheImpl) SetMaxCachedTags(number int) {
	c.logger.Debugf("Set max cached tags to %d", number)

	if number > 0 {
		c.maxTags = number
	}
}

// GetTag Get a tag from cache based on the tag ID.
func (c *CacheImpl) GetTag(tagId string) (*types.IdTagInfo, error) {
	logInfo := c.logger.WithField("tagId", tagId)
	logInfo.Info("Getting a tag from cache")

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
		logInfo.WithError(err).Errorf("Error getting tag from cache")
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

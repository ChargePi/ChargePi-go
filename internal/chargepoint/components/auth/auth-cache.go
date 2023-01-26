package auth

import (
	"encoding/json"
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
	}
)

func NewAuthCache(db *badger.DB) *CacheImpl {
	return &CacheImpl{
		db:      db,
		maxTags: 0,
	}
}

func getTagPrefix(tagId string) []byte {
	return []byte(fmt.Sprintf("cached-tag-%s", tagId))
}

// AddTag Add a tag to the authorization cache.
func (c *CacheImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) {
	logInfo := log.WithField("tagId", tagId)
	logInfo.Debug("Adding a tag to cache")

	// Add a tag if it doesn't exist in the cache already
	err := c.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(getTagPrefix(tagId))
		if err != badger.ErrKeyNotFound {
			return err
		}

		authTag := getTag(tagId, tagInfo)
		if authTag != nil {
			return nil
		}

		err = txn.Set(getTagPrefix(tagId), authTag)
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
	logInfo := log.WithField("tagId", tagId)
	logInfo.Debug("Removing a tag from cache")

	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getTagPrefix(tagId))
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
	logInfo := log.WithField("tagId", tagId)
	logInfo.Info("Getting a tag from cache")

	var tagInfo localauth.AuthorizationData

	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getTagPrefix(tagId))
		if err != badger.ErrKeyNotFound {
			return err
		}

		b, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &tagInfo)
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

package auth

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
)

var (
	ErrTagLimitReached = errors.New("tag limit reached")
)

type (
	LocalAuthList interface {
		AddTag(tagId string, tagInfo *types.IdTagInfo) error
		UpdateTag(tagId string, tagInfo *types.IdTagInfo) error
		RemoveTag(tagId string) error
		RemoveAll()
		GetTag(tagId string) (*types.IdTagInfo, error)
		GetTags() []localauth.AuthorizationData
		SetMaxTags(number int)
		GetVersion() int
		SetVersion(version int)
	}

	LocalAuthListImpl struct {
		db      *badger.DB
		numTags int
		maxTags int
	}
)

func NewLocalAuthList(db *badger.DB, maxTags int) *LocalAuthListImpl {
	return &LocalAuthListImpl{
		db:      db,
		numTags: 0,
		maxTags: maxTags,
	}
}

// AddTag Add a tag to the global authorization cache.
func (l *LocalAuthListImpl) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	if l.numTags+1 >= l.maxTags {
		return ErrTagLimitReached
	}

	return l.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(database.GetLocalAuthTagPrefix(tagId))
		if err != badger.ErrKeyNotFound {
			return err
		}

		authTag := getTag(tagId, tagInfo)
		if authTag != nil {
			return nil
		}

		err = txn.Set(database.GetLocalAuthTagPrefix(tagId), authTag)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

// RemoveTag Remove a tag with the ID from the Local Auth List.
func (l *LocalAuthListImpl) RemoveTag(tagId string) error {
	logInfo := log.WithField("tagId", tagId)
	logInfo.Debug("Removing a tag from local auth list")

	return l.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(database.GetLocalAuthTagPrefix(tagId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

// RemoveAll Remove all tags.
func (l *LocalAuthListImpl) RemoveAll() {
	log.Debugf("Removing local auth list")

	// Remove all cached keys from database
	err := l.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := database.GetLocalAuthTagPrefix("")
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
		log.WithError(err).Error("Error removing local auth list")
	}
}

// GetTag Get a tag
func (l *LocalAuthListImpl) GetTag(tagId string) (*types.IdTagInfo, error) {
	logInfo := log.WithField("tag", tagId)
	logInfo.Info("Fetching the tag")

	var tagInfo localauth.AuthorizationData
	err := l.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(database.GetLocalAuthTagPrefix(tagId))
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
		logInfo.WithError(err).Error("Error fetching local auth tags")
		return nil, err
	}

	return tagInfo.IdTagInfo, nil
}

// GetTags Get all tags stored in the Local Auth store.
func (l *LocalAuthListImpl) GetTags() []localauth.AuthorizationData {
	log.Infof("Fetching tags")
	var tags []localauth.AuthorizationData

	err := l.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := database.GetLocalAuthTagPrefix("")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var data localauth.AuthorizationData
			item := it.Item()

			// Value should be the AuthorizationData struct.
			err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &data)
			})
			if err != nil {
				continue
			}
		}

		return txn.Commit()
	})
	if err != nil {
		log.WithError(err).Error("Error fetching local auth tags")
	}

	return tags
}

// UpdateTag Update a tag in the Local Auth store.
func (l *LocalAuthListImpl) UpdateTag(tagId string, tagInfo *types.IdTagInfo) error {
	logInfo := log.WithField("tagId", tagId)
	logInfo.Info("Updating tag")

	return l.db.Update(func(txn *badger.Txn) error {
		// todo
		return txn.Commit()
	})
}

// GetVersion Get the current version of the Local Auth list.
func (l *LocalAuthListImpl) GetVersion() int {
	version := -1

	err := l.db.View(func(txn *badger.Txn) error {
		versionKey := database.GetLocalAuthVersion()
		item, err := txn.Get(versionKey)
		if err != nil {
			return err
		}

		valueCopy, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		v := binary.BigEndian.Uint32(valueCopy)
		version = int(v)

		return txn.Commit()
	})
	if err != nil {
		return -1
	}

	return version
}

// SetVersion Set the current version of the Local Auth list.
func (l *LocalAuthListImpl) SetVersion(version int) {
	logInfo := log.WithField("version", version)
	logInfo.Info("Updating list version")

	err := l.db.Update(func(txn *badger.Txn) error {
		versionKey := database.GetLocalAuthVersion()

		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, uint32(version))

		err := txn.Set(versionKey, bs)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
	if err != nil {
		logInfo.WithError(err).Error("Error updating list version")
	}
}

// SetMaxTags Set the maximum number of tags that can be stored in the Local Auth list.
func (l *LocalAuthListImpl) SetMaxTags(number int) {
	if number > 0 {
		l.maxTags = number
	}
}

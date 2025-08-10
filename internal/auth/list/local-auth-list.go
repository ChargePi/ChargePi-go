package list

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/ChargePi/ChargePi-go/internal/pkg/database"
	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"go.uber.org/zap"
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

	BadgerLocalAuthList struct {
		db      *badger.DB
		numTags int
		maxTags int
		logger  *zap.Logger
	}
)

func NewLocalAuthList(logger *zap.Logger, db *badger.DB, maxTags int) *BadgerLocalAuthList {
	return &BadgerLocalAuthList{
		db:      db,
		numTags: 0,
		maxTags: maxTags,
		logger:  logger.Named("local-auth-list"),
	}
}

// AddTag Add a tag to the global authorization cache.
func (l *BadgerLocalAuthList) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	logger := l.logger.With(zap.String("tagId", tagId))
	logger.Debug("Adding a tag to local auth list")

	if l.numTags+1 >= l.maxTags {
		return ErrTagLimitReached
	}

	return l.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(database.GetLocalAuthTagPrefix(tagId))
		if !errors.Is(err, badger.ErrKeyNotFound) {
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
func (l *BadgerLocalAuthList) RemoveTag(tagId string) error {
	logger := l.logger.With(zap.String("tagId", tagId))
	logger.Debug("Removing a tag from local auth list")

	return l.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(database.GetLocalAuthTagPrefix(tagId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

// RemoveAll Remove all tags.
func (l *BadgerLocalAuthList) RemoveAll() {
	l.logger.Debug("Removing local auth list")

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
		l.logger.With(zap.Error(err)).Error("Error removing local auth list")
	}
}

// GetTag Get a tag
func (l *BadgerLocalAuthList) GetTag(tagId string) (*types.IdTagInfo, error) {
	logger := l.logger.With(zap.String("tagId", tagId))
	logger.Info("Fetching the tag")

	var tagInfo localauth.AuthorizationData
	err := l.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(database.GetLocalAuthTagPrefix(tagId))
		if !errors.Is(err, badger.ErrKeyNotFound) {
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
		logger.With(zap.Error(err)).Error("Error fetching local auth tags")
		return nil, err
	}

	return tagInfo.IdTagInfo, nil
}

// GetTags Get all tags stored in the Local Auth store.
func (l *BadgerLocalAuthList) GetTags() []localauth.AuthorizationData {
	l.logger.Info("Fetching tags")
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
		l.logger.With(zap.Error(err)).Error("Error fetching local auth tags")
	}

	return tags
}

// UpdateTag Update a tag in the Local Auth store.
func (l *BadgerLocalAuthList) UpdateTag(tagId string, tagInfo *types.IdTagInfo) error {
	logger := l.logger.With(zap.String("tagId", tagId))
	logger.Info("Updating tag")

	return l.db.Update(func(txn *badger.Txn) error {
		// todo
		return txn.Commit()
	})
}

// GetVersion Get the current version of the Local Auth list.
func (l *BadgerLocalAuthList) GetVersion() int {
	l.logger.Info("Fetching list version")
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
func (l *BadgerLocalAuthList) SetVersion(version int) {
	logger := l.logger.With(zap.Int("version", version))
	logger.Info("Updating list version")

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
		logger.With(zap.Error(err)).Error("Error updating list version")
	}
}

// SetMaxTags Set the maximum number of tags that can be stored in the Local Auth list.
func (l *BadgerLocalAuthList) SetMaxTags(number int) {
	if number > 0 {
		l.logger.Sugar().Debugf("Set max tags to %d", number)
		l.maxTags = number
	}
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

package badger

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ChargePi/ChargePi-go/internal/auth/list"
	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/localauth"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

const (
	tagKeyPrefix       = "cached-tag"
	localAuthTagPrefix = "auth-tag"
	localAuthVersion   = "auth-version"
)

func getLocalAuthTagPrefix(tagId string) []byte {
	return []byte(fmt.Sprintf("%s-%s", localAuthTagPrefix, tagId))
}

func getLocalAuthVersion() []byte {
	return []byte(localAuthVersion)
}

func getTagKey(tagId string) []byte {
	return []byte(fmt.Sprintf("%s-%s", tagKeyPrefix, tagId))
}

func (db *Database) AddTagToAuthList(tagId string, tagInfo *types.IdTagInfo) error {
	return db.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(getLocalAuthTagPrefix(tagId))
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}

		authTag := getTag(tagId, tagInfo)
		if authTag != nil {
			return nil
		}

		err = txn.Set(getLocalAuthTagPrefix(tagId), authTag)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (db *Database) RemoveAuthListTag(tagId string) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getLocalAuthTagPrefix(tagId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (db *Database) GetLocalAuthListTag(tagId string) (*types.IdTagInfo, error) {
	var tagInfo localauth.AuthorizationData
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getLocalAuthTagPrefix(tagId))
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
		return nil, err
	}

	return tagInfo.IdTagInfo, nil
}

func (db *Database) GetLocalAuthListTags() ([]localauth.AuthorizationData, error) {
	var tags []localauth.AuthorizationData

	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := getLocalAuthTagPrefix("")
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
		return nil, err
	}

	return tags, nil
}

func (db *Database) GetAuthListTagsForVersion(version int) ([]localauth.AuthorizationData, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) RemoveAuthListAllTagsForVersion(version int) error {
	// Remove all cached keys from database
	return db.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := getLocalAuthTagPrefix("")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := txn.Delete(item.Key())
			if err != nil {
				return err
			}
		}
		return txn.Commit()
	})
}

func (db *Database) AddAuthList(list list.LocalAuthListVersion) error {
	return db.db.Update(func(txn *badger.Txn) error {
		ver := []byte{}
		binary.LittleEndian.PutUint32(ver, uint32(list.Version))

		err := txn.Set(getLocalAuthVersion(), ver)
		if err != nil {
			return err
		}

		// Iterate over all tags and add them to the database
		for _, tag := range list.Tags {
			marshal, err := json.Marshal(tag)
			if err != nil {
				continue
			}

			err = txn.Set(getLocalAuthTagPrefix(tag.IdTag), marshal)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *Database) AddTag(tagId string, tagInfo *types.IdTagInfo) error {
	// Add a tag if it doesn't exist in the cache.
	return db.db.Update(func(txn *badger.Txn) error {
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
}

func (db *Database) RemoveTag(tagId string) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getTagKey(tagId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (db *Database) GetTag(tagId string) (*types.IdTagInfo, error) {
	var tagInfo localauth.AuthorizationData

	err := db.db.View(func(txn *badger.Txn) error {
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
		return nil, err
	}

	return tagInfo.IdTagInfo, nil
}

func (db *Database) GetTags() ([]*types.IdTagInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) RemoveAllTags() error {
	// Remove all cached keys from database
	return db.db.Update(func(txn *badger.Txn) error {
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

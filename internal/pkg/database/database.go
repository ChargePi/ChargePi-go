package database

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
)

var (
	db   *badger.DB
	once = sync.Once{}
)

func Get() *badger.DB {
	once.Do(func() {
		// Load/initialize a database for EVSE, tags, users and settings
		badgerDb, err := badger.Open(badger.DefaultOptions(settings.DatabasePath))
		if err != nil {
			log.WithError(err).Panic("Cannot open/create database")
		}

		db = badgerDb
	})

	return db
}

func GetEvseKey(evseId int) string {
	return fmt.Sprintf("evse-%d", evseId)
}

func GetLocalAuthTagPrefix(tagId string) []byte {
	return []byte(fmt.Sprintf("auth-tag-%s", tagId))
}

func GetLocalAuthVersion() []byte {
	return []byte("auth-version")
}

func GetSmartChargingProfile(profileId int) []byte {
	return []byte(fmt.Sprintf("profile-%d", profileId))
}

func GetOcppConfigurationKey() []byte {
	return []byte(fmt.Sprintf("ocpp-configuration"))
}

func GetEvseSettings(db *badger.DB) []settings.EVSE {
	var evseSettings []settings.EVSE

	// Query the database for EVSE settings.
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte("evse-")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var data settings.EVSE
			item := it.Item()

			// Value should be the EVSE struct.
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
		log.WithError(err).Error("Error querying for EVSE settings")
	}

	return evseSettings
}

package database

import (
	"encoding/json"
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
		opts := badger.DefaultOptions(settings.DatabasePath)
		opts.Logger = newLogger()
		opts.NumGoroutines = 3

		// Load/initialize a database for EVSE, tags, users and settings
		badgerDb, err := badger.Open(opts)
		if err != nil {
			log.WithError(err).Panic("Cannot open/create database")
		}

		db = badgerDb

		// Migrate the database to the latest version
		migration(db)
	})

	return db
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
				log.WithError(err).Warnf("Error unmarshalling EVSE settings for %s", item.Key())
				continue
			}

			evseSettings = append(evseSettings, data)
		}

		return txn.Commit()
	})
	if err != nil {
		log.WithError(err).Error("Error querying for EVSE settings")
	}

	return evseSettings
}

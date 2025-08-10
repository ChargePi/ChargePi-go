package database

import (
	"encoding/json"
	"go.uber.org/zap"
	"sync"

	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/dgraph-io/badger/v3"
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
			zap.L().With(zap.Error(err)).Panic("Cannot open/create database")
		}

		db = badgerDb

		// Migrate the database to the latest version
		migration(db)
	})

	return db
}

func GetEvseSettings(db *badger.DB) []settings.EVSE {
	var evseSettings []settings.EVSE
	logger := zap.L()

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
				logger.With(zap.Error(err)).Sugar().Warnf("Error unmarshalling EVSE settings for %s", item.Key())
				continue
			}

			evseSettings = append(evseSettings, data)
		}

		return txn.Commit()
	})
	if err != nil {
		logger.With(zap.Error(err)).Error("Error querying for EVSE settings")
	}

	return evseSettings
}

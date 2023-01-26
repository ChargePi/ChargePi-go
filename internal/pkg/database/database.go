package database

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
)

func Get() *badger.DB {
	// Load/initialize a database for EVSE, tags, users and settings
	db, err := badger.Open(badger.DefaultOptions(settings.DatabasePath))
	if err != nil {
		log.WithError(err).Panic("Cannot open/create database")
	}

	return db
}

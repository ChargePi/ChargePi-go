package badger

import (
	"github.com/ChargePi/ChargePi-go/internal/users/models"
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
)

type Database struct {
	db     *badger.DB
	logger log.FieldLogger
}

func NewBadgerDb(filePath string) (*Database, error) {
	opts := badger.DefaultOptions(filePath)
	opts.Logger = newLogger()
	opts.NumGoroutines = 1

	// Load/initialize a database for EVSE, tags, users and settings
	badgerDb, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Database{
		db:     badgerDb,
		logger: log.StandardLogger().WithField("component", "badger-db"),
	}, nil
}

func (db *Database) Migrate() {
	db.logger.Debug("Migrating database")

	// todo migration framework?
	_ = db.AddUser(models.User{
		Username: "manufacturer",
		Password: "manufacturer",
		Role:     models.Manufacturer,
	})

	_ = db.AddUser(models.User{
		Username: "technician",
		Password: "technician",
		Role:     models.Technician,
	})

	_ = db.AddUser(models.User{
		Username: "observer",
		Password: "observer",
		Role:     models.Observer,
	})
}

func (db *Database) Close() {
	db.logger.Debug("Closing database")
	err := db.db.Close()
	if err != nil {
		db.logger.WithError(err).Error("Error closing database")
	}
}

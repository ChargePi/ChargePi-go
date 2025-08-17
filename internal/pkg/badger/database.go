package badger

import (
	"context"

	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"

	"github.com/ChargePi/ChargePi-go/internal/users/models"
)

type Database struct {
	db     *badger.DB
	logger *zap.Logger
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
		logger: zap.L().Named("badger_db"),
	}, nil
}

func (db *Database) Migrate() {
	db.logger.Debug("Migrating database")

	// todo migration framework?
	_ = db.AddUser(context.Background(), models.User{
		Username: "manufacturer",
		Password: "manufacturer",
		Role:     models.Manufacturer,
	})

	_ = db.AddUser(context.Background(), models.User{
		Username: "technician",
		Password: "technician",
		Role:     models.Technician,
	})

	_ = db.AddUser(context.Background(), models.User{
		Username: "observer",
		Password: "observer",
		Role:     models.Observer,
	})
}

func (db *Database) Close() {
	db.logger.Debug("Closing database")
	err := db.db.Close()
	if err != nil {
		db.logger.With(zap.Error(err)).Error("Error closing database")
	}
}

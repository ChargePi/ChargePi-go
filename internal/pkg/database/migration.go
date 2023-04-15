package database

import (
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	userDatabase "github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/models"
)

// Initialize the database with default settings.
func migration(db *badger.DB) {
	log.Debug("Migrating database")

	userDb := userDatabase.NewUserDb(db)
	_ = userDb.AddUser(models.User{
		Username: "manufacturer",
		Password: "manufacturer",
		Role:     string(models.Manufacturer),
	})

	_ = userDb.AddUser(models.User{
		Username: "technician",
		Password: "technician",
		Role:     string(models.Technician),
	})

	_ = userDb.AddUser(models.User{
		Username: "observer",
		Password: "observer",
		Role:     string(models.Observer),
	})
}

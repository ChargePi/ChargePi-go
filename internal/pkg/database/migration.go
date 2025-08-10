package database

import (
	userDatabase "github.com/ChargePi/ChargePi-go/internal/users/pkg/database"
	"github.com/ChargePi/ChargePi-go/internal/users/pkg/models"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// Initialize the database with default settings.
func migration(db *badger.DB) {
	zap.L().Debug("Migrating database")

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

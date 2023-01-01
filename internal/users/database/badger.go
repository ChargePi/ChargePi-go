package database

import (
	badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/models"
)

const databasePath = "/tmp/badger"

type BadgerDb struct {
	db *badger.DB
}

func NewBadgerDb() *BadgerDb {
	db, err := badger.Open(badger.DefaultOptions(databasePath))
	if err != nil {
		log.WithError(err).Panic("Cannot open database")
	}

	return &BadgerDb{
		db: db,
	}
}

func (b *BadgerDb) GetUser(username string) (*models.User, error) {
	err := b.db.View(func(txn *badger.Txn) error {

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *BadgerDb) GetUsers() []models.User {
	// b.db.NewTransaction(true).
	return nil
}

func (b *BadgerDb) AddUser(user models.User) error {
	// b.db.NewTransaction(true).
	return nil
}

func (b *BadgerDb) UpdateUser(user models.User) (*models.User, error) {
	// b.db.NewTransaction(true).
	return nil, nil
}

func (b *BadgerDb) DeleteUser(username string) error {
	// b.db.NewTransaction(true).
	return nil
}

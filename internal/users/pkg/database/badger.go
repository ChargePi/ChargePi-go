package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/xBlaz3kx/ChargePi-go/internal/users/pkg/models"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserDoesntExist = errors.New("user doesn't exist")
)

type UserDbBadger struct {
	db *badger.DB
}

func NewUserDb(db *badger.DB) *UserDbBadger {
	return &UserDbBadger{
		db: db,
	}
}

func getUserKey(username string) []byte {
	return []byte(fmt.Sprintf("user-%s", username))
}

func (b *UserDbBadger) GetUser(username string) (*models.User, error) {
	var user models.User
	err := b.db.View(func(txn *badger.Txn) error {
		get, err := txn.Get(getUserKey(username))
		if err != nil {
			return err
		}

		err = get.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (b *UserDbBadger) GetUsers() []models.User {
	var users []models.User

	err := b.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte("user-")
		// Go through every key with “user” prefix.
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var user models.User
			item := it.Item()
			// Value should be the User struct.
			err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &user)
			})
			if err != nil {
				continue
			}
		}

		return txn.Commit()
	})
	if err != nil {
		return nil
	}

	return users
}

func (b *UserDbBadger) AddUser(user models.User) error {
	return b.db.Update(func(txn *badger.Txn) error {
		// Check if user already exists
		_, err := txn.Get(getUserKey(user.Username))
		if err == nil {
			return ErrUserExists
		}

		// todo hash & salt passwords
		marshal, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = txn.Set(getUserKey(user.Username), marshal)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (b *UserDbBadger) UpdateUser(user models.User) (*models.User, error) {
	err := b.db.Update(func(txn *badger.Txn) error {
		// Check if user already exists
		_, err := txn.Get(getUserKey(user.Username))
		if err != nil {
			return ErrUserDoesntExist
		}

		marshal, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = txn.Set(getUserKey(user.Username), marshal)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *UserDbBadger) DeleteUser(username string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getUserKey(username))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

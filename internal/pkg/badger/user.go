package badger

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ChargePi/ChargePi-go/internal/users/models"
	"github.com/dgraph-io/badger/v3"
)

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserDoesntExist = errors.New("user doesn't exist")
)

const userPrefix = "user-"

func (db *Database) GetUser(username string) (*models.User, error) {
	var user models.User
	err := db.db.View(func(txn *badger.Txn) error {
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

func (db *Database) GetUsers() ([]models.User, error) {
	var users []models.User

	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(userPrefix)
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
		return nil, err
	}

	return users, nil
}

func (db *Database) AddUser(user models.User) error {
	return db.db.Update(func(txn *badger.Txn) error {
		// Check if user with username already exists
		_, err := txn.Get(getUserKey(user.Username))
		if err == nil {
			return ErrUserExists
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
}

func (db *Database) UpdateUser(user models.User) (*models.User, error) {
	err := db.db.Update(func(txn *badger.Txn) error {
		// Check if user with username already exists
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

func (db *Database) DeleteUser(username string) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getUserKey(username))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func getUserKey(username string) []byte {
	return []byte(fmt.Sprintf("%s-%s", userPrefix, username))
}

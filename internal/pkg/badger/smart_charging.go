package badger

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

const smartChargingSchedulePrefix = "charging-schedule-"

func getSmartChargingProfile(profileId int) []byte {
	return []byte(fmt.Sprintf("%s-%d", smartChargingSchedulePrefix, profileId))
}

func (db *Database) AddProfile(profile *types.ChargingProfile) error {
	return db.db.Update(func(txn *badger.Txn) error {
		marshal, err := json.Marshal(profile)
		if err != nil {
			return err
		}

		err = txn.Set(getSmartChargingProfile(profile.ChargingProfileId), marshal)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (db *Database) GetProfile(profileId int) (*types.ChargingProfile, error) {
	var profile types.ChargingProfile

	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(getSmartChargingProfile(profileId))
		if err != nil {
			return err
		}

		err = item.Value(func(v []byte) error {
			return json.Unmarshal(v, &profile)
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (db *Database) RemoveProfile(profileId int) error {
	return db.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(getSmartChargingProfile(profileId))
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (db *Database) GetProfiles() ([]types.ChargingProfile, error) {
	var profiles []types.ChargingProfile

	// Query the database for charging profiles
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(smartChargingSchedulePrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var data types.ChargingProfile
			item := it.Item()

			err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &data)
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

	return profiles, nil
}

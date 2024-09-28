package badger

import (
	"encoding/json"
	"fmt"

	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/pkg/ocpp"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ocppConfigurationPrefix = "ocpp-configuration"
	settingsPrefix          = "charge-point-settings"
	evseSettingsPrefix      = "evse-settings"
)

func getOcppConfigurationKey(version ocpp.ProtocolVersion) []byte {
	return []byte(fmt.Sprintf("%s-%s", ocppConfigurationPrefix, version))
}

func getSettingsKey() []byte {
	return []byte(settingsPrefix)
}

func getEvseKey(evseId int) string {
	return fmt.Sprintf("%s-%d", evseSettingsPrefix, evseId)
}

func (db *Database) SetEvseSettings(settings []evse.Settings) error {
	// Sync the settings to the database
	return db.db.Update(func(txn *badger.Txn) error {
		for _, connector := range settings {
			marshal, err := json.Marshal(connector)
			if err != nil {
				return err
			}

			err = txn.Set([]byte(getEvseKey(connector.EvseId)), marshal)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (db *Database) GetEvseSettings() ([]evse.Settings, error) {
	var evseSettings []evse.Settings

	// Query the database for EVSE settings.
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte("evse-")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			var data evse.Settings
			item := it.Item()

			// Value should be the EVSE struct.
			err := item.Value(func(v []byte) error {
				return json.Unmarshal(v, &data)
			})
			if err != nil {
				log.WithError(err).Warnf("Error unmarshalling EVSE settings for %s", item.Key())
				continue
			}

			evseSettings = append(evseSettings, data)
		}

		return txn.Commit()
	})
	if err != nil {
		return nil, errors.Wrap(err, "error getting EVSE settings")
	}

	return evseSettings, nil
}

func (db *Database) GetSettings() (*chargePoint.Settings, error) {
	db.logger.Debug("Getting global settings")

	var settingsS chargePoint.Settings

	err := db.db.View(func(txn *badger.Txn) error {
		config, err := txn.Get(getSettingsKey())
		if err != nil {
			return err
		}

		return config.Value(func(val []byte) error {
			return json.Unmarshal(val, &settingsS)
		})
	})
	if err != nil {
		return nil, err
	}

	return &settingsS, nil
}

func (db *Database) UpdateSettings(settings chargePoint.Settings) error {
	// Read the configuration from the database
	return db.db.Update(func(txn *badger.Txn) error {
		res, err := json.Marshal(settings)
		if err != nil {
			return err
		}

		return txn.Set(getSettingsKey(), res)
	})
}

func (db *Database) GetOcppConfiguration(version int) (*ocpp_v16.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) GeLatestOcppConfiguration() (*ocpp_v16.Config, error) {
	//TODO implement me
	panic("implement me")
}

func (db *Database) SetOcppConfiguration(config ocpp_v16.Config) error {
	//TODO implement me
	panic("implement me")
}

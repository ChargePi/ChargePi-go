package settings

import (
	"encoding/binary"
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var importer Importer

func GetImporter() Importer {
	if importer == nil {
		log.Debug("Creating an importer")
		importer = &ImporterImpl{
			db:              database.Get(),
			settingsManager: GetManager(),
		}
	}

	return importer
}

type Importer interface {
	ImportEVSESettings(settings []settings.EVSE) error
	ImportOcppConfiguration(version configuration.ProtocolVersion, config configuration.Config) error
	ImportLocalAuthList(list settings.AuthList) error
	ImportChargePointSettings(point settings.Settings) error
}

type ImporterImpl struct {
	db              *badger.DB
	settingsManager Manager
}

func (i *ImporterImpl) ImportEVSESettings(settings []settings.EVSE) error {
	log.Debug("Importing connectors to the database")

	// Validate the EVSE settings
	validationErr := validator.New().Struct(settings)
	if validationErr != nil {
		return validationErr
	}

	// Sync the settings to the database
	return i.db.Update(func(txn *badger.Txn) error {
		for _, connector := range settings {
			marshal, err := json.Marshal(connector)
			if err != nil {
				return err
			}

			err = txn.Set([]byte(database.GetEvseKey(connector.EvseId)), marshal)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (i *ImporterImpl) ImportOcppConfiguration(version configuration.ProtocolVersion, config configuration.Config) error {
	log.Debug("Importing ocpp configuration to the database")

	// Validate the settings
	validationErr := validator.New().Struct(config)
	if validationErr != nil {
		return validationErr
	}

	// Sync the settings to the database
	return i.db.Update(func(txn *badger.Txn) error {
		marshal, err := json.Marshal(config)
		if err != nil {
			return err
		}

		return txn.Set(database.GetOcppConfigurationKey(version), marshal)
	})
}

func (i *ImporterImpl) ImportLocalAuthList(list settings.AuthList) error {
	log.Debug("Importing local auth list to the database")

	// Validate the settings
	validationErr := validator.New().Struct(list)
	if validationErr != nil {
		return validationErr
	}

	// Sync the settings to the database
	return i.db.Update(func(txn *badger.Txn) error {
		ver := []byte{}
		binary.LittleEndian.PutUint32(ver, uint32(list.Version))

		err := txn.Set(database.GetLocalAuthVersion(), ver)
		if err != nil {
			return err
		}

		// Iterate over all tags and add them to the database
		for _, tag := range list.Tags {
			marshal, err := json.Marshal(tag)
			if err != nil {
				continue
			}

			err = txn.Set(database.GetLocalAuthTagPrefix(tag.IdTag), marshal)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (i *ImporterImpl) ImportChargePointSettings(settings settings.Settings) error {
	log.Debug("Importing charge point settings to the database")
	return i.settingsManager.SetSettings(settings)
}

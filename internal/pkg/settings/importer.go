package settings

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var importer Importer

func GetImporter() Importer {
	if importer == nil {
		log.Debug("Creating EVSE manager")
		importer = &ImporterImpl{db: database.Get(), ocppVariableManager: ocppConfigManager.GetManager()}
	}

	return importer
}

type Importer interface {
	ImportEVSESettings(settings []settings.EVSE) error
	ImportOcppConfiguration(config configuration.Config) error
	ImportLocalAuthList(list settings.AuthList) error
}

type ImporterImpl struct {
	db                  *badger.DB
	ocppVariableManager ocppConfigManager.Manager
}

func (i *ImporterImpl) ImportEVSESettings(settings []settings.EVSE) error {
	log.Info("Importing connectors to the database")
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

		return txn.Commit()
	})
}

func (i *ImporterImpl) ImportOcppConfiguration(config configuration.Config) error {
	log.Info("Importing ocpp configuration")
	// Sync the settings to the database
	return i.db.Update(func(txn *badger.Txn) error {
		marshal, err := json.Marshal(config)
		if err != nil {
			return err
		}

		err = txn.Set(database.GetOcppConfigurationKey(), marshal)
		if err != nil {
			return err
		}

		return txn.Commit()
	})
}

func (i *ImporterImpl) ImportLocalAuthList(list settings.AuthList) error {
	log.Info("Importing local auth list")
	// Sync the settings to the database
	return i.db.Update(func(txn *badger.Txn) error {

		// txn.Set(database.GetLocalAuthVersion(), nil)

		for _, tag := range list.Tags {
			marshal, err := json.Marshal(tag)
			if err != nil {
				return err
			}

			err = txn.Set(database.GetLocalAuthTagPrefix(tag.IdTag), marshal)
			if err != nil {
				return err
			}
		}

		return txn.Commit()
	})
}

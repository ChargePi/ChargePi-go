package settings

import (
	"encoding/json"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"

	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var (
	manager Manager
	once    = sync.Once{}
)

func init() {
	once.Do(func() {
		GetManager()
		GetExporter()
		GetImporter()
	})
}

func GetManager() Manager {
	if manager == nil {
		log.Debug("Creating settings manager")
		manager = NewManager(database.Get())
	}

	return manager
}

type (
	Manager interface {
		SetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles ...string)
		GetOcppConfiguration(version configuration.ProtocolVersion) ([]core.ConfigurationKey, error)
		GetOcppConfigurationWithKey(version configuration.ProtocolVersion, key string) (*core.ConfigurationKey, error)
		GetChargePointSettings() settings.ChargePoint
		SetChargePointSettings(settings settings.ChargePoint) error
		SetSettings(settings settings.Settings) error
		GetSettings() (*settings.Settings, error)
	}

	Impl struct {
		db                  *badger.DB
		ocppVariableManager ocppConfigManager.Manager
	}
)

func NewManager(db *badger.DB) *Impl {
	return &Impl{
		db:                  db,
		ocppVariableManager: ocppConfigManager.GetManager(),
	}
}

func (i *Impl) SetupOcppConfiguration(version configuration.ProtocolVersion, supportedProfiles ...string) {
	logInfo := log.WithField("version", version)
	logInfo.Info("Setting up OCPP configuration")

	i.ocppVariableManager.SetVersion(version)
	i.ocppVariableManager.SetSupportedProfiles(supportedProfiles...)

	ocppConfig, err := loadConfiguration(i.db, i.ocppVariableManager, version)
	if err != nil {
		logInfo.WithError(err).Errorf("Error loading the configuration from the database")
		return
	}

	err = i.ocppVariableManager.SetConfiguration(*ocppConfig)
	if err != nil {
		logInfo.WithError(err).Errorf("Error setting the configuration to the manager")
		return
	}
}

func (i *Impl) GetChargePointSettings() settings.ChargePoint {
	var settingsS settings.ChargePoint

	err := i.db.View(func(txn *badger.Txn) error {
		config, err := txn.Get(database.GetSettingsKey())
		if err != nil {
			return err
		}

		return config.Value(func(val []byte) error {
			return json.Unmarshal(val, &settingsS)
		})
	})
	if err != nil {
		return settingsS
	}

	return settingsS
}

func (i *Impl) SetChargePointSettings(settings settings.ChargePoint) error {
	// Validate the settings
	validationErr := validator.New().Struct(settings)
	if validationErr != nil {
		return validationErr
	}

	return nil
}

func (i *Impl) SetSettings(settings settings.Settings) error {
	// Validate the settings
	validationErr := validator.New().Struct(settings)
	if validationErr != nil {
		return validationErr
	}

	// Read the configuration from the database
	err := i.db.Update(func(txn *badger.Txn) error {
		res, err := json.Marshal(settings)
		if err != nil {
			return err
		}

		return txn.Set(database.GetSettingsKey(), res)
	})
	return err
}

func (i *Impl) GetSettings() (*settings.Settings, error) {
	var settingsS settings.Settings

	err := i.db.View(func(txn *badger.Txn) error {
		config, err := txn.Get(database.GetSettingsKey())
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

func (i *Impl) GetOcppConfiguration(version configuration.ProtocolVersion) ([]core.ConfigurationKey, error) {
	config, err := loadConfiguration(i.db, i.ocppVariableManager, version)
	if err != nil {
		return nil, err
	}

	err = i.ocppVariableManager.SetConfiguration(*config)
	if err != nil {
		log.WithError(err).Errorf("Error setting the configuration to the manager")
	}

	switch version {
	case configuration.OCPP16:
		return i.ocppVariableManager.GetConfiguration()
	default:
		return nil, nil
	}
}

func (i *Impl) GetOcppConfigurationWithKey(version configuration.ProtocolVersion, key string) (*core.ConfigurationKey, error) {
	switch version {
	case configuration.OCPP16:
		value, err := i.ocppVariableManager.GetConfigurationValue(key)
		if err != nil {
			return nil, err
		}

		return &core.ConfigurationKey{
			Key:      key,
			Readonly: false,
			Value:    value,
		}, nil
	default:
		return nil, nil
	}
}

func loadConfiguration(db *badger.DB, ocppVariableManager ocppConfigManager.Manager, version configuration.ProtocolVersion) (*configuration.Config, error) {
	var ocppConfig configuration.Config

	// Read the configuration from the database
	err := db.View(func(txn *badger.Txn) error {

		config, err := txn.Get(database.GetOcppConfigurationKey(version))
		if err != nil {
			return err
		}

		return config.Value(func(val []byte) error {
			return json.Unmarshal(val, &ocppConfig)
		})
	})
	if err != nil {
		return nil, err
	}

	err = ocppVariableManager.SetConfiguration(ocppConfig)
	if err != nil {
		return nil, err
	}

	return &ocppConfig, nil
}

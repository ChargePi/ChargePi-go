package settings

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dgraph-io/badger/v3"
	"github.com/go-playground/validator/v10"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/ocpp"
	"github.com/xBlaz3kx/ocppManager-go/ocpp_v16"
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
		GetOcppV16Manager() ocpp_v16.Manager
		SetOcppV16Manager(manager ocpp_v16.Manager) error
		GetChargePointSettings() settings.ChargePoint
		SetChargePointSettings(settings settings.ChargePoint) error
		SetSettings(settings settings.Settings) error
		GetSettings() (*settings.Settings, error)
	}

	Impl struct {
		db                    *badger.DB
		ocpp16VariableManager ocpp_v16.Manager
		logger                log.FieldLogger
	}
)

func (i *Impl) GetOcppV16Manager() ocpp_v16.Manager {
	return i.ocpp16VariableManager
}

func (i *Impl) SetOcppV16Manager(manager ocpp_v16.Manager) error {
	if util.IsNilInterfaceOrPointer(manager) {
		return fmt.Errorf("manager cannot be nil")
	}

	i.ocpp16VariableManager = manager
	return nil
}

func NewManager(db *badger.DB) *Impl {
	return &Impl{
		db:     db,
		logger: log.WithField("component", "settings-manager"),
	}
}

func (i *Impl) GetChargePointSettings() settings.ChargePoint {
	i.logger.Debug("Getting charge point settings")
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
	i.logger.Debug("Setting charge point settings")

	// Validate the settings
	validationErr := validator.New().Struct(settings)
	if validationErr != nil {
		return validationErr
	}

	// todo me
	return nil
}

func (i *Impl) SetSettings(settings settings.Settings) error {
	i.logger.Debug("Applying global settings")

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
	i.logger.Debug("Getting global settings")

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

func (i *Impl) GetOcppConfiguration(version ocpp.ProtocolVersion) ([]core.ConfigurationKey, error) {
	i.logger.WithField("version", version).Debug("Getting OCPP configuration")

	switch version {
	case ocpp.OCPP16:
		config, err := loadConfiguration(i.db, i.ocpp16VariableManager)
		if err != nil {
			return nil, err
		}

		err = i.ocpp16VariableManager.SetConfiguration(*config)
		if err != nil {
			i.logger.WithError(err).Errorf("Error setting the configuration to the manager")
		}

		return i.ocpp16VariableManager.GetConfiguration()
	default:
		return nil, fmt.Errorf("unsupported OCPP version: %v", version)
	}
}

func (i *Impl) GetOcppConfigurationWithKey(version ocpp.ProtocolVersion, key string) (*core.ConfigurationKey, error) {
	i.logger.WithField("version", version).Debug("Getting OCPP configuration with key")

	switch version {
	case ocpp.OCPP16:
		value, err := i.ocpp16VariableManager.GetConfigurationValue(ocpp_v16.Key(key))
		if err != nil {
			return nil, err
		}

		return &core.ConfigurationKey{
			Key:      key,
			Readonly: false,
			Value:    value,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported OCPP version: %v", version)
	}
}

func loadConfiguration(db *badger.DB, ocppVariableManager ocpp_v16.Manager) (*ocpp_v16.Config, error) {
	var ocppConfig ocpp_v16.Config

	// Read the configuration from the database
	err := db.View(func(txn *badger.Txn) error {

		config, err := txn.Get(database.GetOcppConfigurationKey(ocpp.OCPP16))
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

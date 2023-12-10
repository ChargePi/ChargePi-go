package settings

import (
	"encoding/binary"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
			logger:          log.WithField("component", "importer"),
		}
	}

	return importer
}

type Importer interface {
	ImportEVSESettings(settings []settings.EVSE) error
	ImportEVSESettingsFromPath(path string) error
	ImportOcppConfiguration(version configuration.ProtocolVersion, config configuration.Config) error
	ImportOcppConfigurationFromPath(version configuration.ProtocolVersion, path string) error
	ImportLocalAuthList(list settings.AuthList) error
	ImportLocalAuthListFromPath(path string) error
	ImportChargePointSettings(point settings.Settings) error
	ImportChargePointSettingsFromPath(path string) error
}

type ImporterImpl struct {
	db              *badger.DB
	settingsManager Manager
	logger          log.FieldLogger
}

func (i *ImporterImpl) ImportEVSESettings(settings []settings.EVSE) error {
	i.logger.Debug("Importing connectors to the database")

	for _, setting := range settings {
		// Validate the EVSE settings
		validationErr := validator.New().Struct(setting)
		if validationErr != nil {
			return validationErr
		}
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
	i.logger.Debug("Importing ocpp configuration to the database")

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
	i.logger.Debug("Importing local auth list to the database")

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
	i.logger.Debug("Importing charge point settings to the database")
	return i.settingsManager.SetSettings(settings)
}

func (i *ImporterImpl) ImportEVSESettingsFromPath(path string) error {
	i.logger.Infof("Importing EVSE settings from %s", path)

	var evseSettings []settings.EVSE

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var s settings.EVSE
			config := viper.New()

			config.SetConfigFile(path)
			err := config.ReadInConfig()
			if err != nil {
				return err
			}

			err = config.Unmarshal(&s)
			if err != nil {
				return err
			}

			evseSettings = append(evseSettings, s)
		}

		return nil
	})
	if err != nil {
		return err
	}

	// Store the settings
	return i.ImportEVSESettings(evseSettings)
}

func (i *ImporterImpl) ImportLocalAuthListFromPath(path string) error {
	i.logger.Infof("Importing tags from %s", path)

	var tagList settings.AuthList
	config := viper.New()

	err := readConfiguration(config, "authList", "yaml", path)
	if err != nil {
		return err
	}

	err = config.Unmarshal(&tagList)
	if err != nil {
		return err
	}

	return i.ImportLocalAuthList(tagList)
}

func (i *ImporterImpl) ImportChargePointSettingsFromPath(path string) error {
	i.logger.Infof("Importing settings from %s", path)

	var cpSettings settings.Settings

	config := viper.New()

	// Read the settings from the file.
	err := readConfiguration(config, "settings", "yaml", path)
	if err != nil {
		return err
	}

	err = config.Unmarshal(&cpSettings)
	if err != nil {
		return err
	}

	return i.ImportChargePointSettings(cpSettings)
}

func (i *ImporterImpl) ImportOcppConfigurationFromPath(version configuration.ProtocolVersion, path string) error {
	i.logger.Infof("Importing OCPP configuration from %s", path)

	var ocppConfiguration configuration.Config
	config := viper.New()

	// Read the settings from the file.
	err := readConfiguration(config, "ocpp", "yaml", path)
	if err != nil {
		return err
	}

	err = config.Unmarshal(&ocppConfiguration)
	if err != nil {
		return err
	}

	return importer.ImportOcppConfiguration(version, ocppConfiguration)
}

func readConfiguration(viper *viper.Viper, fileName, extension, filePath string) error {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}

	return viper.ReadInConfig()
}

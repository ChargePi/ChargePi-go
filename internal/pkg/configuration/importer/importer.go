package importer

import (
	"os"
	"path/filepath"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/auth/list"
	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	"github.com/ChargePi/ChargePi-go/internal/pkg/configuration"
	manager2 "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/pkg/ocpp"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Importer interface {
	ImportEVSESettings(settings []evse.Settings) error
	ImportEVSESettingsFromPath(path string) error
	ImportOcppConfiguration(version ocpp.ProtocolVersion, config ocpp_v16.Config) error
	ImportOcppConfigurationFromPath(version ocpp.ProtocolVersion, path string) error
	ImportLocalAuthList(list list.LocalAuthListVersion) error
	ImportLocalAuthListFromPath(path string) error
	ImportChargePointSettings(point chargepoint.Settings) error
	ImportChargePointSettingsFromPath(path string) error
}

type ImporterImpl struct {
	settingsManager         manager2.Manager
	viper                   *viper.Viper
	evseSettingsRepository  manager.EvseSettingsRepository
	localAuthListRepository auth.LocalAuthListRepository
	logger                  log.FieldLogger
}

func NewImporter(
	settingsManager manager2.Manager,
	evseSettingsRepository manager.EvseSettingsRepository,
	localAuthListRepository auth.LocalAuthListRepository,
) *ImporterImpl {
	return &ImporterImpl{
		settingsManager:         settingsManager,
		viper:                   viper.New(),
		evseSettingsRepository:  evseSettingsRepository,
		localAuthListRepository: localAuthListRepository,
		logger:                  log.StandardLogger().WithField("component", "importer"),
	}
}

func (i *ImporterImpl) ImportEVSESettings(settings []evse.Settings) error {
	i.logger.Debug("Importing connectors to the database")

	for _, setting := range settings {
		// Validate the EVSE settings
		validationErr := validator.New().Struct(setting)
		if validationErr != nil {
			return validationErr
		}
	}

	// Sync the settings to the database
	return i.evseSettingsRepository.SetEvseSettings(settings)
}

func (i *ImporterImpl) ImportOcppConfiguration(version ocpp.ProtocolVersion, config ocpp_v16.Config) error {
	i.logger.Debug("Importing ocpp configuration to the database")

	// Validate the settings
	validationErr := validator.New().Struct(config)
	if validationErr != nil {
		return validationErr
	}

	return i.settingsManager.SetConfiguration(config)
}

func (i *ImporterImpl) ImportLocalAuthList(list list.LocalAuthListVersion) error {
	i.logger.Debug("Importing local auth list to the database")

	// Validate the settings
	validationErr := validator.New().Struct(list)
	if validationErr != nil {
		return validationErr
	}

	return nil

	// Sync the settings to the database
	//return i.db.Update(func(txn *badger.Txn) error {
	//	ver := []byte{}
	//	binary.LittleEndian.PutUint32(ver, uint32(list.Version))
	//
	//	err := txn.Set(badger.getLocalAuthVersion(), ver)
	//	if err != nil {
	//		return err
	//	}
	//
	//	// Iterate over all tags and add them to the database
	//	for _, tag := range list.Tags {
	//		marshal, err := json.Marshal(tag)
	//		if err != nil {
	//			continue
	//		}
	//
	//		err = txn.Set(badger.getLocalAuthTagPrefix(tag.IdTag), marshal)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	return nil
	//})
}

func (i *ImporterImpl) ImportChargePointSettings(settings chargepoint.Settings) error {
	i.logger.Debug("Importing charge point settings to the database")

	// Validate the EVSE settings
	validationErr := validator.New().Struct(settings)
	if validationErr != nil {
		return validationErr
	}

	return i.settingsManager.SetChargePointSettings(settings)
}

func (i *ImporterImpl) ImportEVSESettingsFromPath(path string) error {
	i.logger.Infof("Importing EVSE settings from %s", path)

	var evseSettings []evse.Settings

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var s evse.Settings

		i.viper.SetConfigFile(path)
		err = i.viper.ReadInConfig()
		if err != nil {
			return err
		}

		err = i.viper.Unmarshal(&s)
		if err != nil {
			return err
		}

		evseSettings = append(evseSettings, s)
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

	var tagList list.LocalAuthListVersion

	err := configuration.ReadConfiguration(i.viper, "authList", "yaml", path)
	if err != nil {
		return err
	}

	err = i.viper.Unmarshal(&tagList)
	if err != nil {
		return err
	}

	return i.ImportLocalAuthList(tagList)
}

func (i *ImporterImpl) ImportChargePointSettingsFromPath(path string) error {
	i.logger.Infof("Importing settings from %s", path)

	var cpSettings chargepoint.Settings

	// Read the settings from the file.
	err := configuration.ReadConfiguration(i.viper, "settings", "yaml", path)
	if err != nil {
		return err
	}

	err = i.viper.Unmarshal(&cpSettings)
	if err != nil {
		return err
	}

	return i.ImportChargePointSettings(cpSettings)
}

func (i *ImporterImpl) ImportOcppConfigurationFromPath(version ocpp.ProtocolVersion, path string) error {
	i.logger.Infof("Importing OCPP configuration from %s", path)

	// Todo detect which version of OCPP is being imported from the settings?
	switch version {
	// To

	}
	var ocppConfiguration ocpp_v16.Config
	config := viper.New()

	// Read the settings from the file.
	err := configuration.ReadConfiguration(config, "ocpp", "yaml", path)
	if err != nil {
		return err
	}

	err = config.Unmarshal(&ocppConfiguration)
	if err != nil {
		return err
	}

	return i.ImportOcppConfiguration(version, ocppConfiguration)
}

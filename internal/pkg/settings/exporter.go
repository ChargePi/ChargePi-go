package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/auth"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/database"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var exporter Exporter

func GetExporter() Exporter {
	if exporter == nil {
		log.Debug("Creating an exporter")
		db := database.Get()
		exporter = &ExporterImpl{
			db:              db,
			tagManager:      auth.NewTagManager(db),
			settingsManager: GetManager(),
			logger:          log.StandardLogger().WithField("component", "settings-exporter"),
		}
	}

	return exporter
}

type Exporter interface {
	ExportEVSESettings() []settings.EVSE
	ExportEVSESettingsToFile(path string) error
	ExportOcppConfiguration() configuration.Config
	ExportOcppConfigurationToFile(path string) error
	ExportLocalAuthList() (*settings.AuthList, error)
	ExportLocalAuthListToFile(path string) error
	ExportChargePointSettings() (*settings.Settings, error)
	ExportChargePointSettingsToFile(path string) error
}

type ExporterImpl struct {
	db              *badger.DB
	tagManager      auth.TagManager
	settingsManager Manager
	logger          log.FieldLogger
}

func (i *ExporterImpl) ExportEVSESettings() []settings.EVSE {
	i.logger.Debug("Exporting EVSE settings from the database")
	return database.GetEvseSettings(i.db)
}

func (i *ExporterImpl) ExportOcppConfiguration() configuration.Config {
	i.logger.Debug("Exporting OCPP configuration from the database")
	getConfiguration, _ := i.settingsManager.GetOcppConfiguration(configuration.OCPP16)
	return configuration.Config{
		Version: 1,
		Keys:    getConfiguration,
	}
}

func (i *ExporterImpl) ExportLocalAuthList() (*settings.AuthList, error) {
	i.logger.Debug("Exporting Local auth list from the database")
	return &settings.AuthList{
		Version: i.tagManager.GetAuthListVersion(),
		Tags:    i.tagManager.GetTags(),
	}, nil
}

func (i *ExporterImpl) ExportChargePointSettings() (*settings.Settings, error) {
	i.logger.Debug("Exporting charge point settings from the database")
	return i.settingsManager.GetSettings()
}

func (i *ExporterImpl) ExportEVSESettingsToFile(path string) error {
	i.logger.Infof("Exporting EVSE settings to %s", path)

	evseSettings := i.ExportEVSESettings()
	config := viper.New()

	// Create a file for each EVSE
	for _, evseSetting := range evseSettings {
		fileName := fmt.Sprintf("evse-%d.yaml", evseSetting.EvseId)
		prepareViperCfg(config, fileName, "yaml", path)

		marshal, err := json.Marshal(evseSetting)
		if err != nil {
			return err
		}

		err = config.ReadConfig(bytes.NewBuffer(marshal))
		if err != nil {
			return err
		}

		err = config.WriteConfigAs(filepath.Join(path, fileName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *ExporterImpl) ExportOcppConfigurationToFile(path string) error {
	i.logger.Infof("Exporting OCPP configuration to %s", path)

	ocppConfiguration := exporter.ExportOcppConfiguration()
	config := viper.New()

	prepareViperCfg(config, "ocpp", "yaml", path)

	marshal, err := json.Marshal(ocppConfiguration)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return config.WriteConfigAs(path)
}

func (i *ExporterImpl) ExportLocalAuthListToFile(path string) error {
	i.logger.Infof("Exporting tags to %s", path)

	localAuthList, _ := exporter.ExportLocalAuthList()
	config := viper.New()

	prepareViperCfg(config, "authList", "yaml", path)

	marshal, err := json.Marshal(localAuthList)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return config.WriteConfigAs(path)
}

func (i *ExporterImpl) ExportChargePointSettingsToFile(path string) error {
	i.logger.Infof("Exporting settings to %s", path)

	ocppConfiguration, err := i.ExportChargePointSettings()
	if err != nil {
		return err
	}

	config := viper.New()

	prepareViperCfg(config, "settings", "yaml", path)

	marshal, err := json.Marshal(ocppConfiguration)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return config.WriteConfigAs(path)
}

func prepareViperCfg(viper *viper.Viper, fileName, extension, filePath string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}
}

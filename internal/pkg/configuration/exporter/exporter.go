package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/auth/list"
	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	configManager "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/agrison/go-commons-lang/stringUtils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Exporter interface {
	ExportEVSESettings() ([]evse.Settings, error)
	ExportEVSESettingsToFile(path string) error
	ExportOcppConfiguration() (*ocpp_v16.Config, error)
	ExportOcppConfigurationToFile(path string) error
	ExportLocalAuthList() (*list.LocalAuthListVersion, error)
	ExportLocalAuthListToFile(path string) error
	ExportChargePointSettings() (*chargepoint.Settings, error)
	ExportChargePointSettingsToFile(path string) error
}

type ExporterImpl struct {
	viper                  *viper.Viper
	tagManager             auth.Service
	settingsManager        configManager.Manager
	evseSettingsRepository manager.EvseSettingsRepository
	logger                 log.FieldLogger
}

func NewExporter(tagManager auth.Service, settingsManager configManager.Manager, evseSettingsRepository manager.EvseSettingsRepository) *ExporterImpl {
	return &ExporterImpl{
		viper:                  viper.New(),
		tagManager:             tagManager,
		settingsManager:        settingsManager,
		evseSettingsRepository: evseSettingsRepository,
		logger:                 log.StandardLogger().WithField("component", "settings-exporter"),
	}
}

func (i *ExporterImpl) ExportEVSESettings() ([]evse.Settings, error) {
	i.logger.Debug("Exporting EVSE settings from the database")

	evseSettings, err := i.evseSettingsRepository.GetEvseSettings()
	if err != nil {
		return nil, err
	}

	return evseSettings, nil
}

func (i *ExporterImpl) ExportOcppConfiguration() (*ocpp_v16.Config, error) {
	i.logger.Debug("Exporting OCPP configuration from database")

	getConfiguration, err := i.settingsManager.GetConfiguration()
	if err != nil {
		return nil, err
	}

	return &ocpp_v16.Config{
		Version: 1,
		Keys:    getConfiguration,
	}, nil
}

func (i *ExporterImpl) ExportLocalAuthList() (*list.LocalAuthListVersion, error) {
	i.logger.Debug("Exporting Local auth list from the database")

	tags, err := i.tagManager.GetTags()
	if err != nil {
		return nil, err
	}

	return &list.LocalAuthListVersion{
		Version: i.tagManager.GetAuthListVersion(),
		Tags:    tags,
	}, nil
}

func (i *ExporterImpl) ExportChargePointSettings() (*chargepoint.Settings, error) {
	i.logger.Debug("Exporting charge point settings from the database")

	settings, err := i.settingsManager.GetChargePointSettings()
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (i *ExporterImpl) ExportEVSESettingsToFile(path string) error {
	i.logger.Infof("Exporting EVSE settings to %s", path)

	evseSettings, err := i.ExportEVSESettings()
	if err != nil {
		return err
	}

	// Create a file for each EVSE
	for _, evseSetting := range evseSettings {
		fileName := fmt.Sprintf("evse-%d.yaml", evseSetting.EvseId)
		prepareViperCfg(i.viper, fileName, "yaml", path)

		marshal, err := json.Marshal(evseSetting)
		if err != nil {
			return err
		}

		err = i.viper.ReadConfig(bytes.NewBuffer(marshal))
		if err != nil {
			return err
		}

		err = i.viper.WriteConfigAs(filepath.Join(path, fileName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *ExporterImpl) ExportOcppConfigurationToFile(path string) error {
	i.logger.Infof("Exporting OCPP configuration to %s", path)

	ocppConfiguration, err := i.ExportOcppConfiguration()
	if err != nil {
		return err
	}

	prepareViperCfg(i.viper, "ocpp", "yaml", path)

	marshal, err := json.Marshal(ocppConfiguration)
	if err != nil {
		return err
	}

	err = i.viper.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return i.viper.WriteConfigAs(path)
}

func (i *ExporterImpl) ExportLocalAuthListToFile(path string) error {
	i.logger.Infof("Exporting tags to %s", path)

	localAuthList, err := i.ExportLocalAuthList()
	if err != nil {
		return err
	}

	prepareViperCfg(i.viper, "authList", "yaml", path)
	marshal, err := json.Marshal(localAuthList)
	if err != nil {
		return err
	}

	err = i.viper.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return i.viper.WriteConfigAs(path)
}

func (i *ExporterImpl) ExportChargePointSettingsToFile(path string) error {
	i.logger.Infof("Exporting settings to %s", path)

	ocppConfiguration, err := i.ExportChargePointSettings()
	if err != nil {
		return err
	}

	prepareViperCfg(i.viper, "settings", "yaml", path)

	marshal, err := json.Marshal(ocppConfiguration)
	if err != nil {
		return err
	}

	err = i.viper.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return i.viper.WriteConfigAs(path)
}

func prepareViperCfg(viper *viper.Viper, fileName, extension, filePath string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}
}

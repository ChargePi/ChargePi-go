package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/spf13/viper"
)

func ExportLocalAuthList(exporter Exporter, config *viper.Viper, filePath string) error {
	localAuthList, _ := exporter.ExportLocalAuthList()

	prepareViperCfg(config, "authList", "yaml", filePath)

	marshal, err := json.Marshal(localAuthList)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return config.WriteConfigAs(filePath)
}

func ExportOcppConfiguration(exporter Exporter, config *viper.Viper, filePath string) error {
	ocppConfiguration := exporter.ExportOcppConfiguration()

	prepareViperCfg(config, "ocpp", "yaml", filePath)

	marshal, err := json.Marshal(ocppConfiguration)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return config.WriteConfigAs(filePath)
}

func ExportSettings(exporter Exporter, config *viper.Viper, filePath string) error {
	ocppConfiguration, err := exporter.ExportChargePointSettings()
	if err != nil {
		return err
	}

	prepareViperCfg(config, "settings", "yaml", filePath)

	marshal, err := json.Marshal(ocppConfiguration)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	return config.WriteConfigAs(filePath)
}

func ExportEVSEs(exporter Exporter, config *viper.Viper, path string) error {
	evseSettings := exporter.ExportEVSESettings()

	log.Debug(evseSettings)
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

func prepareViperCfg(viper *viper.Viper, fileName, extension, filePath string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}
}

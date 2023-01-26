package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/spf13/viper"
)

func ExportLocalAuthList(exporter Exporter, config *viper.Viper, filePath string) error {
	localAuthList, _ := exporter.ExportLocalAuthList()

	marshal, err := json.Marshal(localAuthList)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	prepareViperCfg(config, "authList", "yaml", filePath)

	return config.WriteConfigAs(fmt.Sprintf("%s%dauth.yaml", filePath, os.PathListSeparator))
}

func ExportOcppConfiguration(exporter Exporter, config *viper.Viper, filePath string) error {
	ocppConfiguration := exporter.ExportOcppConfiguration()

	marshal, err := json.Marshal(ocppConfiguration)
	if err != nil {
		return err
	}

	err = config.ReadConfig(bytes.NewBuffer(marshal))
	if err != nil {
		return err
	}

	prepareViperCfg(config, "ocpp", "yaml", filePath)

	return config.WriteConfigAs(fmt.Sprintf("%s%docpp.yaml", filePath, os.PathListSeparator))
}

func ExportEVSEs(exporter Exporter, config *viper.Viper, path string) error {
	evseSettings := exporter.ExportEVSESettings()

	// Create a file for each EVSE
	for _, evseSetting := range evseSettings {
		marshal, err := json.Marshal(evseSetting)
		if err != nil {
			return err
		}

		err = config.ReadConfig(bytes.NewBuffer(marshal))
		if err != nil {
			return err
		}

		fileName := fmt.Sprintf("evse-%d", evseSetting.EvseId)
		prepareViperCfg(config, fileName, "yaml", path)

		err = config.WriteConfigAs(fmt.Sprintf("%s%d%s", path, os.PathListSeparator, fileName))
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

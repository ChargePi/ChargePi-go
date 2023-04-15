package settings

import (
	"os"
	"path/filepath"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

func ImportLocalAuthList(importer Importer, config *viper.Viper, filePath string) error {
	var tagList settings.AuthList

	err := readConfiguration(config, "authList", "yaml", filePath)
	if err != nil {
		return err
	}

	err = config.Unmarshal(&tagList)
	if err != nil {
		return err
	}

	return importer.ImportLocalAuthList(tagList)
}

func ImportOcppConfiguration(importer Importer, config *viper.Viper, filePath, version string) error {
	var ocppConfiguration configuration.Config

	// Read the settings from the file.
	err := readConfiguration(config, "ocpp", "yaml", filePath)
	if err != nil {
		return err
	}

	err = config.Unmarshal(&ocppConfiguration)
	if err != nil {
		return err
	}

	return importer.ImportOcppConfiguration(configuration.ProtocolVersion(version), ocppConfiguration)
}

func ImportSettings(importer Importer, config *viper.Viper, filePath string) error {
	var cpSettings settings.Settings

	// Read the settings from the file.
	err := readConfiguration(config, "settings", "yaml", filePath)
	if err != nil {
		return err
	}

	err = config.Unmarshal(&cpSettings)
	if err != nil {
		return err
	}

	return importer.ImportChargePointSettings(cpSettings)
}

// ImportEVSEs loads connectors from a folder and imports them into the database.
func ImportEVSEs(importer Importer, config *viper.Viper, path string) error {
	var evseSettings []settings.EVSE

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var s settings.EVSE

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
	return importer.ImportEVSESettings(evseSettings)
}

func readConfiguration(viper *viper.Viper, fileName, extension, filePath string) error {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}

	return viper.ReadInConfig()
}

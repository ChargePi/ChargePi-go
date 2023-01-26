package settings

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/agrison/go-commons-lang/stringUtils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

func ImportLocalAuthList(importer Importer, config *viper.Viper, filePath string) error {
	var tagList settings.AuthList

	readConfiguration(config, "authList", "yaml", filePath)

	err := config.Unmarshal(&tagList)
	if err != nil {
		return err
	}

	return importer.ImportLocalAuthList(tagList)
}

func ImportOcppConfiguration(importer Importer, config *viper.Viper, filePath string) error {
	var ocppConfiguration configuration.Config

	readConfiguration(config, "ocpp", "yaml", filePath)
	err := config.Unmarshal(&ocppConfiguration)
	if err != nil {
		return err
	}

	return importer.ImportOcppConfiguration(ocppConfiguration)
}

// ImportEVSEs loads a connector from file
func ImportEVSEs(importer Importer, config *viper.Viper, path string) error {
	var evseSettings []settings.EVSE

	for _, paths := range getDirPaths(path) {
		var s []settings.EVSE

		readConfiguration(config, "evse", "yaml", paths)
		err := config.Unmarshal(&s)
		if err != nil {
			return err
		}
	}

	// Store the settings
	return importer.ImportEVSESettings(evseSettings)
}

func readConfiguration(viper *viper.Viper, fileName, extension, filePath string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatalf("Cannot parse config file")
	}

	log.Debugf("Using configuration file: %s", viper.ConfigFileUsed())
}

func getDirPaths(path string) []string {
	paths := []string{path}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return paths
	}

	if fileInfo.IsDir() {
		paths = []string{}
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return paths
		}

		for _, file := range files {
			paths = append(paths, fmt.Sprintf("%s%d%s", path, os.PathSeparator, file.Name()))
		}
	}

	return paths
}

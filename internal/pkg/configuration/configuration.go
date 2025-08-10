package configuration

import (
	"encoding/json"
	"go.uber.org/zap"
	"strings"

	"github.com/ChargePi/ChargePi-go/internal/pkg/settings"
	"github.com/go-playground/validator/v10"

	settingsModel "github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/spf13/viper"
)

func InitSettings(settingsFilePath string) {
	config := viper.GetViper()
	setupEnv(config)
	setDefaults(config)
	readConfiguration(config, "settings", "yaml", settingsFilePath)
}

func readConfiguration(viper *viper.Viper, fileName, extension, filePath string) {
	logger := zap.L()
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)
	viper.AddConfigPath(settingsModel.CurrentFolder)
	viper.AddConfigPath(settingsModel.EvseFolder)
	viper.AddConfigPath(settingsModel.DockerFolder)

	// If a file path is provided, use that instead of the default ones.
	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}

	// Read the configuration file.
	err := viper.ReadInConfig()
	switch err {
	case nil:
		logger.Sugar().Debugf("Using configuration file: %s", viper.ConfigFileUsed())
	default:
		logger.With(zap.Error(err)).Warn("Cannot parse config file")
	}
}

func setupEnv(viper *viper.Viper) {
	viper.SetEnvPrefix("chargepi")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func setDefaults(viper *viper.Viper) {
	viper.SetDefault(settingsModel.Model, "ChargePi")
	viper.SetDefault(settingsModel.Vendor, "xBlaz3kx")
	viper.SetDefault(settingsModel.ProtocolVersion, "1.6")
}

// GetSettings gets settings from cache or reads the settings file if the cached settings are not found.
func GetSettings() *settingsModel.Settings {
	logger := zap.L()
	logger.Info("Fetching settings..")

	var conf settingsModel.Settings

	if stringUtils.IsEmpty(viper.ConfigFileUsed()) {
		logger.Debug("Using database settings..")

		// Load the settings persisted in the database.
		getSettings, settingsErr := settings.GetManager().GetSettings()
		if settingsErr != nil {
			logger.With(zap.Error(settingsErr)).Fatal("Cannot load settings from database")
		}

		marshal, settingsErr := json.Marshal(getSettings)
		if settingsErr != nil {
			return nil
		}

		// Load the settings into viper.
		settingsErr = viper.ReadConfig(strings.NewReader(string(marshal)))
		if settingsErr != nil {
			return nil
		}
	}

	err := viper.Unmarshal(&conf)
	if err != nil {
		logger.With(zap.Error(err)).Fatal("Cannot unmarshal settings")
	}

	// Validate the settings
	validationErr := validator.New().Struct(conf)
	if validationErr != nil {
		logger.With(zap.Error(validationErr)).Fatal("Invalid settings")
	}

	return &conf
}

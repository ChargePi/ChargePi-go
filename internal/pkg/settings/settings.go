package settings

import (
	"strings"
	"sync"

	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
)

var (
	EVSESettings = sync.Map{}
)

func InitSettings(settingsFilePath string) {
	setupEnv()
	setDefaults()
	readConfiguration(viper.GetViper(), "settings", "yaml", settingsFilePath)
}

func readConfiguration(viper *viper.Viper, fileName, extension, filePath string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)
	viper.AddConfigPath(settings.CurrentFolder)
	viper.AddConfigPath(settings.EvseFolder)
	viper.AddConfigPath(settings.DockerFolder)

	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatalf("Cannot parse config file")
	}

	log.Debugf("Using configuration file: %s", viper.ConfigFileUsed())
}

func setupEnv() {
	viper.SetEnvPrefix("chargepi")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func setDefaults() {
	viper.SetDefault(settings.Model, "ChargePi")
	viper.SetDefault(settings.Vendor, "xBlaz3kx")
	viper.SetDefault(settings.MaxChargingTime, 180)
	viper.SetDefault(settings.ProtocolVersion, "1.6")
}

func SetupOcppConfiguration(db *badger.DB, version configuration.ProtocolVersion, supportedProfiles ...string) {
	// todo load the settings from database if they're stored. Else use defaults.
	ocppConfigManager.SetVersion(version)
	ocppConfigManager.SetSupportedProfiles(supportedProfiles...)

	err := ocppConfigManager.GetManager().SetConfiguration(configuration.Config{})
	if err != nil {
		return
	}
}

// GetSettings gets settings from cache or reads the settings file if the cached settings are not found.
func GetSettings() *settings.Settings {
	log.Debug("Fetching settings..")
	// todo load settings from the database if they're presisted there.
	// overwrite if they're already set from viper

	var conf settings.Settings

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.WithError(err).Fatalf("Couldn't load settings")
	}

	validationErr := validator.New().Struct(conf)
	if validationErr != nil {
		log.WithError(validationErr).Fatalf("Invalid settings")
	}

	return &conf
}

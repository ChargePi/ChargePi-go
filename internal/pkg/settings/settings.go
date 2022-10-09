package settings

import (
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/go-playground/validator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/util"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"path/filepath"
	"strings"
	"sync"
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
	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	} else {
		viper.SetConfigName(fileName)
		viper.SetConfigType(extension)
		viper.AddConfigPath(currentFolder)
		viper.AddConfigPath(evseFolder)
		viper.AddConfigPath(dockerFolder)
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
	viper.SetDefault(Model, "ChargePi")
	viper.SetDefault(Vendor, "xBlaz3kx")
	viper.SetDefault(MaxChargingTime, 180)
	viper.SetDefault(ProtocolVersion, "1.6")
	viper.SetDefault(LoggingFormat, "gelf")
}

func SetupOcppConfigurationManager(filePath string, version configuration.ProtocolVersion, supportedProfiles ...string) {
	fileName := strings.TrimSuffix(filePath, filepath.Ext(filePath))

	ocppConfigManager.SetFileFormat(util.JSON)
	ocppConfigManager.SetVersion(version)
	ocppConfigManager.SetFileName(filepath.Base(fileName))
	ocppConfigManager.SetFilePath(filepath.Dir(filePath))
	ocppConfigManager.SetSupportedProfiles(supportedProfiles...)

	// Load the configuration
	err := ocppConfigManager.LoadConfiguration()
	if err != nil {
		log.WithError(err).Fatalf("Cannot load OCPP configuration")
	}
}

// GetSettings gets settings from cache or reads the settings file if the cached settings are not found.
func GetSettings() *settings.Settings {
	log.Debug("Fetching settings..")
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

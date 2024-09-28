package configuration

import (
	"strings"

	"github.com/ChargePi/ChargePi-go/pkg/ocpp"
	"github.com/agrison/go-commons-lang/stringUtils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DatabasePath  = "/tmp/chargepi"
	CurrentFolder = "./configs"
	EvseFolder    = "./configs/evses"
	DockerFolder  = "/etc/ChargePi/configs"
)

// Configuration variables
const (
	Model           = "chargepoint.info.ocpp.model"
	Vendor          = "chargepoint.info.ocpp.vendor"
	MaxChargingTime = "chargepoint.info.maxChargingTime"
	ProtocolVersion = "chargepoint.info.protocolVersion"

	Debug      = "debug"
	ApiEnabled = "api.enabled"
	ApiAddress = "api.address"
	ApiPort    = "api.port"
)

// Flags
const (
	DebugFlag          = "debug"
	ApiAddressFlag     = "api-address"
	SettingsFlag       = "settings"
	EvseFlag           = "evse"
	AuthFileFlag       = "auth"
	OcppConfigPathFlag = "ocpp"
	OcppVersion        = "v"
)

var DefaultVendor = "xBlaz3kx"
var DefaultModel = "ChargePi"

func InitSettings(settingsFilePath string) {
	config := viper.GetViper()

	setupEnv(config)
	setDefaults(config)

	err := ReadConfiguration(config, "settings", "yaml", settingsFilePath)
	if err != nil {
		log.WithError(err).Fatalf("Cannot read configuration file")
	}
}

func ReadConfiguration(v *viper.Viper, fileName, extension, filePath string) error {
	v.SetConfigName(fileName)
	v.SetConfigType(extension)
	v.AddConfigPath(CurrentFolder)
	v.AddConfigPath(EvseFolder)
	v.AddConfigPath(DockerFolder)

	// If a file path is provided, use that instead of the default ones.
	if stringUtils.IsNotEmpty(filePath) {
		v.SetConfigFile(filePath)
	}

	// Read the configuration file.
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("No configuration file found")
			return nil
		}
	}

	return err
}

func setupEnv(viper *viper.Viper) {
	viper.SetEnvPrefix("chargepi")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}

func setDefaults(viper *viper.Viper) {
	viper.SetDefault(Model, DefaultModel)
	viper.SetDefault(Vendor, DefaultVendor)
	viper.SetDefault(ProtocolVersion, ocpp.OCPP16)
}

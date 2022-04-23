package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/agrison/go-commons-lang/stringUtils"
	"github.com/go-playground/validator"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	currentFolder   = "./configs"
	connectorFolder = "./configs/connectors"
	dockerFolder    = "/etc/ChargePi/configs"

	Model           = "chargepoint.info.ocpp.model"
	Vendor          = "chargepoint.info.ocpp.vendor"
	MaxChargingTime = "chargepoint.info.maxChargingTime"
	ProtocolVersion = "chargepoint.info.protocolVersion"
	LoggingFormat   = "chargepoint.logging.format"
	Debug           = "debug"
	ApiEnabled      = "api.enabled"
	ApiAddress      = "api.address"
	ApiPort         = "api.port"
)

var (
	ConnectorSettings = sync.Map{}
)

func InitSettings(settingsFilePath string) {
	readConfiguration(viper.GetViper(), "settings", "yaml", settingsFilePath)
	setupEnv()
	setDefaults()
}

func readConfiguration(viper *viper.Viper, fileName, extension, filePath string) {
	if stringUtils.IsNotEmpty(filePath) {
		viper.SetConfigFile(filePath)
	} else {
		viper.SetConfigName(fileName)
		viper.SetConfigType(extension)
		viper.AddConfigPath(currentFolder)
		viper.AddConfigPath(connectorFolder)
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

	ocppConfigManager.SetFileFormat(JSON)
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

// loadConnectorFromPath loads a connector from file
func loadConnectorFromPath(name, path string) (*settings.Connector, error) {
	// Read the connector settings from the file in the directory
	var (
		connectorCfg = viper.New()
		connector    settings.Connector
	)

	readConfiguration(connectorCfg, name, "json", path)

	err := connectorCfg.Unmarshal(&connector)
	if err != nil {
		log.WithError(err).Errorf("Cannot read connector file")
		return nil, err
	}

	log.Debugf("Read connector from %s", path)
	cachePathKey := fmt.Sprintf("connectorEvse%dId%d", connector.EvseId, connector.ConnectorId)
	ConnectorSettings.Store(cachePathKey, &connector)

	return &connector, nil
}

// GetConnectors Scan the connectors folder, read all the connectors' settings and cache the settings.
func GetConnectors(connectorsFolderPath string) []*settings.Connector {
	var (
		connectors []*settings.Connector
	)

	log.Debug("Fetching connectors..")

	err := filepath.Walk(connectorsFolderPath, func(path string, info os.FileInfo, err error) error {
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Load a connector from the path
		connector, err := loadConnectorFromPath(info.Name(), path)
		if err != nil {
			return err
		}

		// Append the configuration to the array
		connectors = append(connectors, connector)
		return nil
	})

	if err != nil {
		log.WithError(err).Errorf("Error reading connectors")
	}

	return connectors
}

// UpdateConnectorStatus update the Connector's status in the connector configuration file
func UpdateConnectorStatus(evseId, connectorId int, status core.ChargePointStatus) {
	var (
		cachePathKey = fmt.Sprintf("connectorEvse%dId%d", evseId, connectorId)
		connector    settings.Connector
		err          error
		logInfo      = log.WithFields(log.Fields{
			"evseId":      evseId,
			"connectorId": connectorId,
			"status":      status,
		})
	)

	viperCfg, isFound := ConnectorSettings.Load(cachePathKey)
	if !isFound {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	cfg := viperCfg.(*viper.Viper)

	err = cfg.Unmarshal(&connector)
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	// Replace the connector status
	cfg.Set("status", status)

	err = cfg.WriteConfig()
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	logInfo.Debugf("Updated status at connector %d", connectorId)
}

// UpdateConnectorSessionInfo update the Connector's Session object in the connector configuration file
func UpdateConnectorSessionInfo(evseId, connectorId int, session *settings.Session) {
	var (
		cachePathKey = fmt.Sprintf("connectorEvse%dId%d", evseId, connectorId)
		connector    *settings.Connector
		err          error
		logInfo      = log.WithFields(log.Fields{
			"evseId":      evseId,
			"connectorId": connectorId,
			"session":     session,
		})
	)

	logInfo.Debugf("Updating session info")
	viperCfg, isFound := ConnectorSettings.Load(cachePathKey)
	if !isFound {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	cfg := viperCfg.(*viper.Viper)

	err = cfg.Unmarshal(&connector)
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	connector.Session = *session

	marshal, err := json.Marshal(connector)
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	err = cfg.ReadConfig(bytes.NewReader(marshal))
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	err = cfg.WriteConfig()
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	logInfo.Debugf("Updated session for connector")
}

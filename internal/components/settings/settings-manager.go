package settings

import (
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	goCache "github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/cache"
	ocppConfigManager "github.com/xBlaz3kx/ocppManager-go"
	"github.com/xBlaz3kx/ocppManager-go/configuration"
	"os"
	"path/filepath"
	"strings"
)

const (
	Key = "settings"

	currentFolder = "./connectors"
	dockerFolder  = "/etc/ChargePi/configs"
)

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

// loadSettings from file
func loadSettings(cache *goCache.Cache, path string) *settings.Settings {
	log.Debug("Fetching settings..")
	var conf settings.Settings

	err := fig.Load(&conf,
		fig.File(filepath.Base(path)),
		fig.Dirs(filepath.Dir(path), dockerFolder),
	)
	if err != nil {
		log.WithError(err).Fatalf("Couldn't load settings")
	}

	log.Debugf("Fetched settings from %s", path)
	cache.Set(Key, &conf, goCache.NoExpiration)

	return &conf
}

// GetSettings gets settings from cache or reads the settings file if the cached settings are not found.
func GetSettings(cache *goCache.Cache, settingsPath string) *settings.Settings {
	cacheSettings, isFound := cache.Get(Key)
	if isFound {
		return cacheSettings.(*settings.Settings)
	}

	return loadSettings(cache, settingsPath)
}

// loadConnectorFromPath loads a connector from file
func loadConnectorFromPath(cache *goCache.Cache, name, path string) (*settings.Connector, error) {
	// Read the connector settings from the file in the directory
	var connector settings.Connector
	err := fig.Load(&connector,
		fig.File(name),
		fig.Dirs(dockerFolder+"/connectors", currentFolder, filepath.Dir(path)),
	)
	if err != nil {
		log.WithError(err).Errorf("Cannot read connector file")
		return nil, err
	}

	log.Tracef("Read connector from %s", path)

	if cache != nil {
		// Add the connector config file path to cache
		err = cache.Add(fmt.Sprintf("connectorEvse%dId%dFilePath", connector.EvseId, connector.ConnectorId), path, goCache.NoExpiration)
		if err != nil {
			return nil, err
		}

		// Add the Connector configuration to the cache
		err = cache.Add(fmt.Sprintf("connectorEvse%dId%dConfiguration", connector.EvseId, connector.ConnectorId), &connector, goCache.NoExpiration)
		if err != nil {
			return nil, err
		}
	}

	return &connector, nil
}

// GetConnectors Scan the connectors folder and read all the connectors' settings.
func GetConnectors(cache *goCache.Cache, connectorsFolderPath string) []*settings.Connector {
	var (
		connectors []*settings.Connector
	)

	err := filepath.Walk(connectorsFolderPath, func(path string, info os.FileInfo, err error) error {
		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Load a connector from the path
		connector, err := loadConnectorFromPath(cache, info.Name(), path)
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
		cachePathKey      = fmt.Sprintf("connectorEvse%dId%dFilePath", evseId, connectorId)
		connectorSettings settings.Connector
		err               error
		logInfo           = log.WithFields(log.Fields{
			"evseId":      evseId,
			"connectorId": connectorId,
			"status":      status,
		})
	)

	// Get the file path from cache
	result, isFound := cache.Cache.Get(cachePathKey)
	if !isFound {
		logInfo.Debugf("Path of the file not found in cache")
		return
	}

	connectorFilePath := result.(string)
	err = fig.Load(&connectorSettings,
		fig.File(filepath.Base(connectorFilePath)),
		fig.Dirs(filepath.Dir(connectorFilePath)))
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	// Replace the connector status
	connectorSettings.Status = string(status)
	err = WriteToFile(connectorFilePath, &connectorSettings)
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating connector status")
		return
	}

	logInfo.Debugf("Updated status at connector %d", connectorId)
}

// UpdateConnectorSessionInfo update the Connector's Session object in the connector configuration file
func UpdateConnectorSessionInfo(evseId, connectorId int, session *settings.Session) {
	var (
		cachePathKey      = fmt.Sprintf("connectorEvse%dId%dFilePath", evseId, connectorId)
		cacheConnectorKey = fmt.Sprintf("connectorEvse%dId%dConfiguration", evseId, connectorId)
		connectorSettings *settings.Connector
		err               error
		logInfo           = log.WithFields(log.Fields{
			"evseId":      evseId,
			"connectorId": connectorId,
			"session":     session,
		})
	)

	logInfo.Debugf("Updating session info")
	// Get the file path from cache
	result, isFound := cache.Cache.Get(cachePathKey)
	if !isFound {
		logInfo.Errorf("Path of the file not found in cache")
		return
	}
	var connectorFilePath = result.(string)

	// Try to find the connector's settings in the cache, if it fails, get settings from the file
	cachedSettings, isFound := cache.Cache.Get(cacheConnectorKey)
	if isFound {
		connectorSettings = cachedSettings.(*settings.Connector)
	} else {
		err = fig.Load(&connectorSettings,
			fig.File(filepath.Base(connectorFilePath)),
			fig.Dirs(filepath.Dir(connectorFilePath)))
		if err != nil {
			logInfo.WithError(err).Errorf("Error updating session info")
			return
		}
	}

	// Replace the session values
	connectorSettings.Session = *session

	err = WriteToFile(connectorFilePath, &connectorSettings)
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating session info")
		return
	}

	logInfo.Debugf("Updated session for connector")
}

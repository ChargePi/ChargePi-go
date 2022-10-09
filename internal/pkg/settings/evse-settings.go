package settings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/models/settings"
	"os"
	"path/filepath"
)

// loadEVSEFromPath loads a connector from file
func loadEVSEFromPath(name, path string) (*settings.EVSE, error) {
	// Read the evse settings from the file in the directory
	var (
		evseCfg = viper.New()
		evse    settings.EVSE
	)

	readConfiguration(evseCfg, name, "json", path)

	err := evseCfg.Unmarshal(&evse)
	if err != nil {
		log.WithError(err).Errorf("Cannot read evse file")
		return nil, err
	}

	log.Debugf("Loaded evse from %s", path)
	cachePathKey := fmt.Sprintf("evse%d", evse.EvseId)
	EVSESettings.Store(cachePathKey, &evse)

	return &evse, nil
}

// GetEVSEs Scan the evse folder, read all the evses' settings and cache the settings.
func GetEVSEs(evseFolderPath string) []*settings.EVSE {
	var (
		evses []*settings.EVSE
	)

	log.Debug("Fetching evses..")
	err := filepath.Walk(evseFolderPath, func(path string, info os.FileInfo, err error) error {
		// Skip (sub) directories
		if info.IsDir() {
			return nil
		}

		// Load an evse from the path
		evse, err := loadEVSEFromPath(info.Name(), path)
		if err != nil {
			return err
		}

		// Append the configuration to the array
		evses = append(evses, evse)
		return nil
	})

	if err != nil {
		log.WithError(err).Errorf("Error reading evses")
	}

	return evses
}

// UpdateEVSEStatus update the Connector's status in the connector configuration file
func UpdateEVSEStatus(evseId int, status core.ChargePointStatus) {
	var (
		cachePathKey = fmt.Sprintf("evse%d", evseId)
		evse         settings.EVSE
		err          error
		logInfo      = log.WithFields(log.Fields{
			"evseId": evseId,
			"status": status,
		})
	)

	viperCfg, isFound := EVSESettings.Load(cachePathKey)
	if !isFound {
		logInfo.WithError(err).Errorf("Error loading evse configuration")
		return
	}

	cfg := viperCfg.(*viper.Viper)

	err = cfg.Unmarshal(&evse)
	if err != nil {
		logInfo.WithError(err).Errorf("Error unmarshalling evse configuration")
		return
	}

	// Replace the evse status
	cfg.Set("status", status)

	err = cfg.WriteConfig()
	if err != nil {
		logInfo.WithError(err).Errorf("Error updating evse status")
		return
	}

	logInfo.Debugf("Updated status at evse %d", evseId)
}

// UpdateSession update the Connector's Session object in the connector configuration file
func UpdateSession(evseId int, session *settings.Session) {
	var (
		cachePathKey = fmt.Sprintf("evse%d", evseId)
		connector    *settings.EVSE
		err          error
		logInfo      = log.WithFields(log.Fields{
			"evseId":  evseId,
			"session": session,
		})
	)

	logInfo.Debugf("Updating session info")
	viperCfg, isFound := EVSESettings.Load(cachePathKey)
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

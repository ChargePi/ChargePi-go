package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kkyr/fig"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GetSettings Read settings from the specified path
func GetSettings() {
	var (
		settings     Settings
		settingsPath = ""
		err          error
	)
	cacheSettings, isFound := cache.Cache.Get("settingsFilePath")
	if !isFound {
		log.Fatal("settings file path not found")
	}
	settingsPath = cacheSettings.(string)
	err = fig.Load(&settings,
		fig.File(filepath.Base(settingsPath)),
		fig.Dirs(filepath.Dir(settingsPath)),
	)
	if err != nil {
		log.Fatal(err)
	}
	cache.Cache.Set("settings", &settings, goCache.NoExpiration)
	log.Println("Read settings from ", settingsPath)
}

// GetConnectors Scan the connectors folder and read all the connectors' settings.
func GetConnectors() []*Connector {
	var (
		connectors           []*Connector
		connectorsFolderPath = ""
	)

	connectorPath, isFound := cache.Cache.Get("connectorsFolderPath")
	if isFound {
		connectorsFolderPath = connectorPath.(string)
	}

	err := filepath.Walk(connectorsFolderPath, func(path string, info os.FileInfo, err error) error {
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Read the connector settings from the file in the directory
		var connector Connector
		err = fig.Load(&connector,
			fig.File(info.Name()),
			fig.Dirs("./connectors", filepath.Dir(path)),
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// Append the configuration to the array
		connectors = append(connectors, &connector)
		log.Println("Read connector from ", path)
		// Add the Connector configuration to the cache
		err = cache.Cache.Add(fmt.Sprintf("connectorEvse%dId%dFilePath", connector.EvseId, connector.ConnectorId), path, goCache.NoExpiration)
		err = cache.Cache.Add(fmt.Sprintf("connectorEvse%dId%dConfiguration", connector.EvseId, connector.ConnectorId), &connector, goCache.NoExpiration)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
	return connectors
}

// UpdateConnectorStatus update the Connector's status in the connector configuration file
func UpdateConnectorStatus(evseId int, connectorId int, status core.ChargePointStatus) {
	var (
		cachePathKey      = fmt.Sprintf("connectorEvse%dId%dFilePath", evseId, connectorId)
		connectorSettings Connector
		err               error
	)
	//Get the file path from cache
	result, isFound := cache.Cache.Get(cachePathKey)
	if !isFound {
		fmt.Println("Path of the file not found in cache")
		return
	}

	connectorFilePath := result.(string)
	err = fig.Load(&connectorSettings,
		fig.File(filepath.Base(connectorFilePath)),
		fig.Dirs(filepath.Dir(connectorFilePath)))
	if err != nil {
		log.Println("Error updating connector status: ", err)
		return
	}

	//Replace the connector status
	connectorSettings.Status = string(status)
	err = WriteToFile(connectorFilePath, &connectorSettings)
	if err != nil {
		log.Println("Error updating connector status: ", err)
		return
	}
	log.Println("Updated status at connector", connectorId)
}

// UpdateConnectorSessionInfo update the Connector's Session object in the connector configuration file
func UpdateConnectorSessionInfo(evseId int, connectorId int, session *Session) {
	var (
		cachePathKey      = fmt.Sprintf("connectorEvse%dId%dFilePath", evseId, connectorId)
		cacheConnectorKey = fmt.Sprintf("connectorEvse%dId%dConfiguration", evseId, connectorId)
		connectorSettings *Connector
		err               error
	)

	log.Println("Updating session info for connector ", connectorId)
	//Get the file path from cache
	result, isFound := cache.Cache.Get(cachePathKey)
	if !isFound {
		fmt.Println("Path of the file not found in cache")
		return
	}
	var connectorFilePath = result.(string)

	// Try to find the connector's settings in the cache, if it fails, get settings from the file
	settings, isFound := cache.Cache.Get(cacheConnectorKey)
	if isFound {
		connectorSettings = settings.(*Connector)
	} else {
		err := fig.Load(&connectorSettings,
			fig.File(filepath.Base(connectorFilePath)),
			fig.Dirs(filepath.Dir(connectorFilePath)))
		if err != nil {
			log.Println("Error updating connector status: ", err)
			return
		}
	}

	//Replace the session values
	connectorSettings.Session = struct {
		IsActive      bool   `fig:"IsActive"`
		TransactionId string `fig:"TransactionId" default:""`
		TagId         string `fig:"TagId" default:""`
		Started       string `fig:"Started" default:""`
		Consumption   []types.MeterValue
	}(*session)

	err = WriteToFile(connectorFilePath, &connectorSettings)
	if err != nil {
		log.Println("Error updating session info: ", err)
		return
	}
	log.Println("Updated session for connector ", connectorId)
}

func WriteToFile(filename string, structure interface{}) (err error) {
	var (
		encodingType string
		marshal      []byte
	)
	// Check if the file format is supported
	splitFile := strings.Split(filename, ".")
	isValidFile, err := regexp.MatchString("^.*\\.(json|yaml|yml)$", filename)
	if len(splitFile) > 0 && isValidFile {
		encodingType = splitFile[len(splitFile)-1]
	}

	switch encodingType {
	case "yaml", "yml":
		marshal, err = yaml.Marshal(&structure)
		break
	case "json":
		marshal, err = json.MarshalIndent(&structure, "", "\t")
		break
	default:
		return errors.New("unsupported file format")
	}

	err = ioutil.WriteFile(filename, marshal, 0644)
	if err != nil {
		log.Println(err)
	}
	return err
}

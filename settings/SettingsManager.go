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

type Settings struct {
	ChargePoint struct {
		Info struct {
			Vendor               string `fig:"Vendor" default:"UL FE"`
			Model                string `fig:"Model" default:"ChargePi"`
			Id                   string `fig:"Id" validate:"required"`
			ProtocolVersion      string `fig:"ProtocolVersion" default:"1.6"`
			CurrentClientVersion string `fig:"CurrentClientVersion" default:"1.0"`
			TargetClientVersion  string `fig:"TargetClientVersion" default:"1.0"`
			ServerUri            string `fig:"ServerUri" validate:"required"`
			LogServer            string `fig:"LogServer" validate:"required"`
			MaxChargingTime      int    `fig:"MaxChargingTime" default:"180"`
			TLS                  struct {
				IsEnabled             bool   `fig:"isEnabled"`
				CACertificatePath     string `fig:"CACertificatePath"`
				ClientCertificatePath string `fig:"ClientCertificatePath"`
				ClientKeyPath         string `fig:"ClientKeyPath"`
			}
		}
		Hardware struct {
			Lcd struct {
				IsSupported bool   `fig:"IsSupported"`
				I2CAddress  string `fig:"I2CAddress"`
			}
			TagReader struct {
				IsSupported bool   `fig:"IsSupported"`
				ReaderModel string `fig:"ReaderModel"`
				Device      string `fig:"Device"`
				ResetPin    int    `fig:"ResetPin"`
			}
			LedIndicator struct {
				Enabled          bool   `fig:"Enabled"`
				DataPin          int    `fig:"DataPin"`
				IndicateCardRead bool   `fig:"IndicateCardRead"`
				Type             string `fig:"Type"`
				Invert           bool   `fig:"Invert"`
			}
			PowerMeters struct {
				MinPower int `fig:"MinPower" default:"20"`
				Retries  int `fig:"Retries" default:"3"`
			}
		}
	}
}

type Connector struct {
	EvseId      int    `fig:"EvseId" validate:"required"`
	ConnectorId int    `fig:"ConnectorId" validate:"required"`
	Type        string `fig:"Type" validate:"required"`
	Status      string `fig:"Status" validation:"required"`
	Session     struct {
		IsActive      bool   `fig:"IsActive"`
		TransactionId string `fig:"TransactionId" default:""`
		TagId         string `fig:"TagId" default:""`
		Started       string `fig:"Started" default:""`
		Consumption   []types.MeterValue
	} `fig:"Session"`
	Relay struct {
		RelayPin     int  `fig:"RelayPin" validate:"required"`
		InverseLogic bool `fig:"InverseLogic"`
	} `fig:"Relay"`
	PowerMeter struct {
		Enabled              bool    `fig:"Enabled"`
		PowerMeterPin        int     `fig:"PowerMeterPin"`
		SpiBus               int     `fig:"SpiBus" default:"0"`
		PowerUnits           string  `fig:"PowerUnits" `
		Consumption          float64 `fig:"Consumption"`
		ShuntOffset          float64 `fig:"ShuntOffset"`
		VoltageDividerOffset float64 `fig:"VoltageDividerOffset"`
	} `fig:"PowerMeter"`
}

type Session struct {
	IsActive      bool
	TransactionId string
	TagId         string
	Started       string
	Consumption   []types.MeterValue
}

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
		var connector Connector
		err = fig.Load(&connector,
			fig.File(info.Name()),
			fig.Dirs("./connectors", filepath.Dir(path)),
		)
		if err != nil {
			fmt.Println(err)
			return err
		}
		connectors = append(connectors, &connector)
		log.Println("Read connector from ", path)
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
	//Replace the session values
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
	//Try to find Connector settings in the cache, if fails, get them from the file
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
	splitFile := strings.Split(filename, ".")
	isValidFile, err := regexp.MatchString("^.*\\.(json|yaml|yml)$", filename)
	if len(splitFile) > 0 && isValidFile {
		encodingType = splitFile[len(splitFile)-1]
	}
	switch encodingType {
	case "yaml":
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

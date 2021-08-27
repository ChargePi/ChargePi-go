package settings

import (
	"encoding/json"
	"errors"
	"fmt"
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
			Vendor               string
			Model                string
			Id                   string
			ProtocolVersion      string
			CurrentClientVersion string
			TargetClientVersion  string
			ServerUri            string
			LogServer            string
			MaxChargingTime      int
			TLS                  struct {
				IsEnabled             bool
				CACertificatePath     string
				ClientCertificatePath string
				ClientKeyPath         string
			}
		}
		Hardware struct {
			Lcd struct {
				IsSupported bool
				I2CAddress  string
			}
			TagReader struct {
				IsSupported bool
				ReaderModel string
				Device      string
				ResetPin    int
			}
			LedIndicator struct {
				Enabled          bool
				DataPin          string
				IndicateCardRead bool
				Type             string
				Invert           bool
			}
			PowerMeters struct {
				MinPower int
				Retries  int
			}
		}
	}
}

type Connector struct {
	EvseId      int
	ConnectorId int
	Type        string
	Status      string
	Session     struct {
		IsActive      bool
		TransactionId string
		TagId         string
		Started       string
		Consumption   []types.MeterValue
	}
	Relay struct {
		RelayPin     int
		InverseLogic bool
	}
	PowerMeter struct {
		Enabled              bool
		PowerMeterPin        int
		SpiBus               int
		PowerUnits           string
		Consumption          float64
		ShuntOffset          float64
		VoltageDividerOffset float64
	} `json:"powerMeter,omitempty" yaml:"powerMeter,omitempty"`
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
	var settings Settings
	var settingsPath = ""
	cacheSettings, isFound := cache.Cache.Get("settingsFilePath")
	if isFound {
		settingsPath = cacheSettings.(string)
	}
	DecodeFile(settingsPath, &settings)
	err := cache.Cache.Add("settings", settings, goCache.NoExpiration)
	if err != nil {
		panic(err)
	}
	log.Println("Added settings to cache")
}

// GetConnectors Scan the connectors folder and read all the connectors' settings.
func GetConnectors() []*Connector {
	var connectors []*Connector
	var connectorsFolderPath = ""
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
		DecodeFile(path, &connector)
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

func UpdateConnectorStatus(evseId int, connectorId int, status core.ChargePointStatus) {
	log.Println("Updating status of the connector ", connectorId)
	var cachePathKey = fmt.Sprintf("connectorEvse%dId%dFilePath", evseId, connectorId)
	var connectorSettings *Connector
	//Get the file path from cache
	result, isFound := cache.Cache.Get(cachePathKey)
	if !isFound {
		fmt.Println("Path of the file not found in cache")
		return
	}
	DecodeFile(result.(string), &connectorSettings)
	//Replace the session values
	connectorSettings.Status = string(status)
	err := WriteToFile(result.(string), &connectorSettings)
	if err != nil {
		log.Println("Error updating connector status: ", err)
		return
	}
	log.Println("Updated status ", connectorId)
}

func UpdateConnectorSessionInfo(evseId int, connectorId int, session *Session) {
	log.Println("Updating session info for connector ", connectorId)
	var cachePathKey = fmt.Sprintf("connectorEvse%dId%dFilePath", evseId, connectorId)
	var cacheConnectorKey = fmt.Sprintf("connectorEvse%dId%dConfiguration", evseId, connectorId)
	var connectorSettings *Connector
	//Get the file path from cache
	result, isFound := cache.Cache.Get(cachePathKey)
	if !isFound {
		fmt.Println("Path of the file not found in cache")
		return
	}
	//Try to find Connector settings in the cache, if fails, get them from the file
	settings, isFound := cache.Cache.Get(cacheConnectorKey)
	if isFound {
		connectorSettings = settings.(*Connector)
	} else {
		DecodeFile(result.(string), &connectorSettings)
	}
	//Replace the session values
	connectorSettings.Session = *session
	err := WriteToFile(result.(string), &connectorSettings)
	if err != nil {
		log.Println("Error updating session info: ", err)
		return
	}
	log.Println("Updated session for ", connectorId)
}

func DecodeFile(filename string, structure interface{}) {
	var (
		err          error
		encodingType string
		file         *os.File
	)
	file, encodingType, err = getFile(filename)
	if err != nil {
		fmt.Println(filename)
		panic(err)
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	switch encodingType {
	case "yaml":
		yaml.Unmarshal(byteValue, &structure)
		break
	case "json":
		json.Unmarshal(byteValue, &structure)
		break
	}
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

func getFile(filename string) (file *os.File, encodingType string, err error) {
	encodingType = "json"
	splitFile := strings.Split(filename, ".")
	isValidFile, err := regexp.MatchString("^.*\\.(json|yaml|yml)$", filename)
	if !isValidFile {
		err = errors.New("invalid file structure")
		return
	}
	if len(splitFile) > 0 {
		encodingType = splitFile[len(splitFile)-1]
	}
	fmt.Println("Checking if file exists: ", filename)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("File not found:", filename)
		return
	}
	file, err = os.Open(filename)
	if err != nil {
		return
	}
	return
}

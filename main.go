package main

import (
	"flag"
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	cache2 "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/data"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func initLogger(grayLogAddress string) {
	gelfWriter, err := gelf.NewWriter(grayLogAddress)
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
}

func getFlags() error {
	var workingDirectory, _ = os.Getwd()
	var configurationPath = fmt.Sprintf("%s/configs", workingDirectory)
	// Get the paths from arguments
	configurationFileFormatFlag := flag.String("config-format", "json", "Format of the configuration files (YAML, JSON or TOML)")
	configurationFileFormat := strings.ToLower(*configurationFileFormatFlag)
	configurationFolderPath := flag.String("config-folder", configurationPath, "Path to the configuration folder")
	connectorsFolderPath := flag.String("connector-folder", fmt.Sprintf("%s/%s", configurationPath, "connectors"), "Path to the connector folder")
	configurationFilePath := flag.String("ocpp-config", fmt.Sprintf("%s/%s.%s", configurationPath, "configuration", configurationFileFormat), "Path to the OCPP configuration file")
	settingsFilePath := flag.String("settings", fmt.Sprintf("%s/%s.%s", configurationPath, "settings", configurationFileFormat), "Path to the settings file")
	authFilePath := flag.String("auth", fmt.Sprintf("%s/%s.%s", configurationPath, "auth", configurationFileFormat), "Path to the authorization persistence file")
	flag.Parse()
	// Put the paths in the Cache
	err := cache.Cache.Add("configurationFolderPath", *configurationFolderPath, cache2.NoExpiration)
	err = cache.Cache.Add("connectorsFolderPath", *connectorsFolderPath, cache2.NoExpiration)
	err = cache.Cache.Add("configurationFilePath", *configurationFilePath, cache2.NoExpiration)
	err = cache.Cache.Add("settingsFilePath", *settingsFilePath, cache2.NoExpiration)
	err = cache.Cache.Add("authFilePath", *authFilePath, cache2.NoExpiration)
	return err
}

func main() {
	var (
		config     *settings.Settings
		tagReader  *hardware.TagReader
		lcd        *hardware.LCD
		ledStrip   *hardware.LEDStrip
		connectors []*settings.Connector
		ledError   error
		handler    chargepoint.ChargePointHandler
	)
	quitChannel := make(chan os.Signal, 1)
	err := getFlags()
	if err != nil {
		return
	}
	// Read settings from files -> SettingsManager
	settings.GetSettings()
	go settings.InitConfiguration()
	go data.GetAuthFile()
	connectors = settings.GetConnectors()
	cacheSettings, isFound := cache.Cache.Get("settings")
	if !isFound {
		panic("settings not found")
	}
	config = cacheSettings.(*settings.Settings)
	initLogger(config.ChargePoint.Info.LogServer)
	if config.ChargePoint.Hardware.TagReader.IsSupported && config.ChargePoint.Hardware.TagReader.ReaderModel == "PN532" {
		// Make a TagReader object if configured
		log.Println("Preparing tag reader from config:", config.ChargePoint.Hardware.TagReader.ReaderModel)
		tagChannel := make(chan string)
		tagReader = &hardware.TagReader{
			TagChannel:       tagChannel,
			DeviceConnection: config.ChargePoint.Hardware.TagReader.Device,
			ResetPin:         config.ChargePoint.Hardware.TagReader.ResetPin,
		}
		// 2. Listen to RFID/NFC reader -> separate thread that communicates with the ChargePoint struct
		go tagReader.ListenForTags()
	}
	if config.ChargePoint.Hardware.Lcd.IsSupported {
		// 3. LCD listens to ChargePoint struct for messages on a separate thread
		log.Println("Preparing LCD from config")
		lcdChannel := make(chan hardware.LCDMessage, 5)
		lcd = hardware.NewLCD(lcdChannel)
		go lcd.DisplayMessages()
	}
	if config.ChargePoint.Hardware.LedIndicator.Enabled == true && config.ChargePoint.Hardware.LedIndicator.Type == "WS281x" {
		// Add an LED strip if configured
		log.Println("Preparing LED strip from config: ", config.ChargePoint.Hardware.LedIndicator.Type)
		stripLength := len(connectors)
		if config.ChargePoint.Hardware.LedIndicator.IndicateCardRead {
			stripLength++
		}
		ledStrip, ledError = hardware.NewLEDStrip(stripLength, config.ChargePoint.Hardware.LedIndicator.DataPin)
		if ledError != nil {
			log.Println(ledError)
		}
	}
	if config.ChargePoint.Info.ProtocolVersion == "1.6" {
		// 4. Instantiate ChargePoint struct with provided settings
		handler := &chargepoint.ChargePointHandler{
			Settings:  config,
			TagReader: tagReader,
			LCD:       lcd,
			LEDStrip:  ledStrip,
		}
		handler.AddConnectors(connectors)
		// Listen to incoming requests from the Central System
		handler.Run()
	} else {
		log.Fatal("Protocol version not supported: ", config.ChargePoint.Info.ProtocolVersion)
	}

	//capture terminate signal
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Received a signal to terminate..")
	handler.CleanUp(core.ReasonPowerLoss)
	fmt.Println("Exiting...")
}

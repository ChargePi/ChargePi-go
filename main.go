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
	"github.com/xBlaz3kx/ChargePi-go/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func initLogger(grayLogAddress string) {
	gelfWriter, err := gelf.NewWriter(grayLogAddress)
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
}

func getFlags() {
	once := sync.Once{}
	once.Do(func() {
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
		cache.Cache.Set("configurationFolderPath", *configurationFolderPath, cache2.NoExpiration)
		cache.Cache.Set("connectorsFolderPath", *connectorsFolderPath, cache2.NoExpiration)
		cache.Cache.Set("configurationFilePath", *configurationFilePath, cache2.NoExpiration)
		cache.Cache.Set("settingsFilePath", *settingsFilePath, cache2.NoExpiration)
		cache.Cache.Set("authFilePath", *authFilePath, cache2.NoExpiration)
	})
}

func readSettings() {
	// Read settings from files
	settings.GetSettings()
	go settings.InitConfiguration()
	go data.GetAuthFile()
}

func main() {
	var (
		config      *settings.Settings
		tagReader   reader.Reader
		lcd         display.LCD
		handler     *chargepoint.ChargePointHandler
		quitChannel = make(chan os.Signal, 1)
	)

	// read flags and settings file
	getFlags()
	readSettings()

	cacheSettings, isFound := cache.Cache.Get("settings")
	if !isFound {
		panic("settings not found")
	}
	config = cacheSettings.(*settings.Settings)
	chargePointInfo := config.ChargePoint.Info

	initLogger(chargePointInfo.LogServer)

	switch settings.ProtocolVersion(chargePointInfo.ProtocolVersion) {
	case settings.OCPP16:
		// try to create hardware components based on settings
		tagReader = reader.NewTagReader()
		lcd = display.NewDisplay()

		if tagReader != nil {
			go tagReader.ListenForTags()
		}

		if lcd != nil {
			go lcd.ListenForMessages()
		}

		// Instantiate ChargePoint struct with provided settings
		handler = &chargepoint.ChargePointHandler{
			Settings:  config,
			TagReader: tagReader,
			LCD:       lcd,
		}

		// Listen to incoming requests from the Central System
		handler.Run()
		break
	case settings.OCPP201:
		log.Println("Version 2.0.1 is not supported yet.")
	default:
		log.Fatal("Protocol version not supported:", chargePointInfo.ProtocolVersion)
	}

	//capture terminate signal
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Received a signal to terminate..")
	handler.CleanUp(core.ReasonPowerLoss)
	fmt.Println("Exiting...")
}

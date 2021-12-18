package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	goCache "github.com/patrickmn/go-cache"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint/v16"
	"github.com/xBlaz3kx/ChargePi-go/components/cache"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/display"
	"github.com/xBlaz3kx/ChargePi-go/components/hardware/reader"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/conf-manager"
	"github.com/xBlaz3kx/ChargePi-go/components/settings/settings-manager"
	"github.com/xBlaz3kx/ChargePi-go/data/auth"
	"github.com/xBlaz3kx/ChargePi-go/data/settings"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

var once = sync.Once{}

func initLogger(grayLogAddress string) {
	gelfWriter, err := gelf.NewWriter(grayLogAddress)
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}

	log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
}

func getFlags() {
	once.Do(func() {
		var (
			workingDirectory, _ = os.Getwd()
			configurationPath   = fmt.Sprintf("%s/configs", workingDirectory)
		)

		// Get the paths from arguments
		configurationFileFormatFlag := flag.String("config-format", "json", "Format of the configuration files (YAML, JSON or TOML)")
		configurationFileFormat := strings.ToLower(*configurationFileFormatFlag)
		configurationFolderPath := flag.String("config-folder", configurationPath, "Path to the configuration folder")
		connectorsFolderPath := flag.String("connector-folder", fmt.Sprintf("%s/%s", configurationPath, "connectors"), "Path to the connector folder")
		configurationFilePath := flag.String("ocpp-config", fmt.Sprintf("%s/%s.%s", configurationPath, "configuration", configurationFileFormat), "Path to the OCPP configuration file")
		settingsFilePath := flag.String("settings", fmt.Sprintf("%s/%s.%s", configurationPath, "settings", configurationFileFormat), "Path to the settings file")
		authFilePath := flag.String("auth", fmt.Sprintf("%s/%s.%s", configurationPath, "auth", configurationFileFormat), "Path to the authorization persistence file")
		flag.Parse()

		// Put the paths in the mem
		cache.Cache.Set("configurationFolderPath", *configurationFolderPath, goCache.NoExpiration)
		cache.Cache.Set("connectorsFolderPath", *connectorsFolderPath, goCache.NoExpiration)
		cache.Cache.Set("configurationFilePath", *configurationFilePath, goCache.NoExpiration)
		cache.Cache.Set("settingsFilePath", *settingsFilePath, goCache.NoExpiration)
		cache.Cache.Set("authFilePath", *authFilePath, goCache.NoExpiration)
	})
}

func readSettings() {
	// Read settings from files
	settings_manager.GetSettings()
	go conf_manager.InitConfiguration()
	go auth.LoadAuthFile()
}

func main() {
	var (
		config      *settings.Settings
		tagReader   reader.Reader
		lcd         display.LCD
		err         error
		handler     *v16.ChargePointHandler
		quitChannel = make(chan os.Signal, 1)
		ctx, cancel = context.WithCancel(context.Background())
	)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	// Read flags and settings file
	getFlags()
	readSettings()

	cacheSettings, isFound := cache.Cache.Get("settings")
	if !isFound {
		log.Fatalf("settings not found")
	}
	config = cacheSettings.(*settings.Settings)

	var (
		chargePointInfo = config.ChargePoint.Info
		hardware        = config.ChargePoint.Hardware
	)

	// Create the logger
	initLogger(chargePointInfo.LogServer)

	// Create hardware components based on settings
	tagReader, err = reader.NewTagReader(hardware.TagReader)
	if err == nil {
		go tagReader.ListenForTags(ctx)
	}

	lcd, err = display.NewDisplay(hardware.Lcd)
	if err == nil {
		go lcd.ListenForMessages(ctx)
	}

	switch settings.ProtocolVersion(chargePointInfo.ProtocolVersion) {
	case settings.OCPP16:
		// Create the client & listen to incoming requests from the Central System
		handler = v16.NewChargePoint(tagReader, lcd)
		handler.Run(ctx, config)
		break
	case settings.OCPP201:
		log.Println("Version 2.0.1 is not supported yet.")
	default:
		log.Fatal("Protocol version not supported:", chargePointInfo.ProtocolVersion)
	}

	// Capture terminate signal
	<-quitChannel
	handler.CleanUp(core.ReasonPowerLoss)
	cancel()
	fmt.Println("Exiting...")
}

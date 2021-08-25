package main

import (
	"fmt"
	"github.com/Graylog2/go-gelf/gelf"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/xBlaz3kx/ChargePi-go/cache"
	"github.com/xBlaz3kx/ChargePi-go/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/hardware"
	"github.com/xBlaz3kx/ChargePi-go/settings"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func initLogger(grayLogAddress string) {
	gelfWriter, err := gelf.NewWriter(grayLogAddress)
	if err != nil {
		log.Fatalf("gelf.NewWriter: %s", err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, gelfWriter))
}

func main() {
	var (
		config     settings.Settings
		tagReader  *hardware.TagReader
		lcd        *hardware.LCD
		ledStrip   *hardware.LEDStrip
		connectors []*settings.Connector
		ledError   error
	)
	// 1. Read settings from files -> SettingsManager
	cacheSettings, isFound := cache.Cache.Get("settings")
	if !isFound {
		panic("settings not found")
	}
	config = cacheSettings.(settings.Settings)
	connectors = settings.GetConnectors()
	initLogger(config.ChargePoint.Info.LogServer)
	quitChannel := make(chan os.Signal, 1)
	if config.ChargePoint.Hardware.TagReader.IsSupported && config.ChargePoint.Hardware.TagReader.ReaderModel == "PN532" {
		log.Println("Preparing tag reader from config:", config.ChargePoint.Hardware.TagReader.ReaderModel)
		tagChannel := make(chan string)
		// 2. Listen to RFID/NFC reader -> separate thread that communicates with the ChargePoint struct
		tagReader = &hardware.TagReader{
			TagChannel:       tagChannel,
			DeviceConnection: config.ChargePoint.Hardware.TagReader.Device,
			ResetPin:         config.ChargePoint.Hardware.TagReader.ResetPin,
		}
		go tagReader.ListenForTags()
	}
	if config.ChargePoint.Hardware.Lcd.IsSupported {
		log.Println("Preparing LCD from config")
		lcdChannel := make(chan hardware.LCDMessage, 5)
		// 3. LCD listens to ChargePoint struct for messages on a separate thread
		lcd = hardware.NewLCD(lcdChannel)
		go lcd.DisplayMessages()
	}
	if config.ChargePoint.Hardware.LedIndicator.Enabled == true && config.ChargePoint.Hardware.LedIndicator.Type == "WS281x" {
		log.Println("Preparing LED strip from config: ", config.ChargePoint.Hardware.LedIndicator.Type)
		dataPin, err := strconv.Atoi(config.ChargePoint.Hardware.LedIndicator.DataPin)
		if err != nil {
			dataPin = 18
		}
		stripLength := len(connectors)
		if config.ChargePoint.Hardware.LedIndicator.IndicateCardRead {
			stripLength++
		}
		ledStrip, ledError = hardware.NewLEDStrip(stripLength, dataPin)
		if ledError != nil {
			log.Println(ledError)
		}
	}
	var handler chargepoint.ChargePointHandler
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
	}
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel
	fmt.Println("Received a signal to terminate..")
	close(quitChannel)
	handler.CleanUp(core.ReasonPowerLoss)
	fmt.Println("Exiting...")
}

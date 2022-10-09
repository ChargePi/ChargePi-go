package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
)

var (
	// Basic configuration settings
	configurationFilePath string
	evseFolderPath        string
	settingsFilePath      string
	authFilePath          string

	rootCmd = &cobra.Command{
		Use:   "chargepi",
		Short: "ChargePi is an open-source OCPP client.",
		Long:  ``,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	settings.InitSettings(settingsFilePath)

	var (
		isDebug      = viper.GetBool(settings.Debug)
		mainSettings = settings.GetSettings()
		evses        = settings.GetEVSEs(evseFolderPath)
	)

	// Run the charge point
	chargepoint.Run(isDebug, mainSettings, evses, configurationFilePath, authFilePath)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.WithError(err).Fatal("Unable to run")
	}
}

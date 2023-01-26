package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	settings2 "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
)

var (
	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the ChargePi core",
		Long:  ``,
		Run:   run,
	}

	// Basic configuration setting
	settingsFilePath string
)

func run(cmd *cobra.Command, args []string) {
	settings.InitSettings(settingsFilePath)

	var (
		debug        = viper.GetBool(settings2.Debug)
		mainSettings = settings.GetSettings()
	)

	// Run the actual charge point
	chargepoint.Run(debug, mainSettings)
}

func init() {
	rootCmd.AddCommand(runCmd)

	rootCmd.PersistentFlags().StringVar(&settingsFilePath, settings2.SettingsFlag, "", "config file path")

	// Here you will define your flags and configuration settings.
	runCmd.PersistentFlags().String(settings2.ApiAddressFlag, "localhost", "address of the api")
	runCmd.PersistentFlags().Int(settings2.ApiPortFlag, 4269, "port for the API")

	_ = viper.BindPFlag(settings2.ApiAddress, runCmd.PersistentFlags().Lookup(settings2.ApiAddressFlag))
	_ = viper.BindPFlag(settings2.ApiPort, runCmd.PersistentFlags().Lookup(settings2.ApiPortFlag))
}

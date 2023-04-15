package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/chargepoint"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/configuration"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
)

var (
	// Basic configuration setting
	settingsFilePath string
)

// runCommand is the command for the ChargePi core.
func runCommand() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the ChargePi core",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			configuration.InitSettings(settingsFilePath)

			debug := viper.GetBool(settings.Debug)
			mainSettings := configuration.GetSettings()

			// Run the charge point
			chargepoint.Run(debug, mainSettings)
		},
	}

	runCmd.Flags().StringVar(&settingsFilePath, settings.SettingsFlag, "", "config file path")
	runCmd.Flags().String(settings.ApiAddressFlag, "localhost:4269", "listen address")

	_ = viper.BindPFlag(settings.ApiAddress, runCmd.Flags().Lookup(settings.ApiAddressFlag))

	return runCmd
}

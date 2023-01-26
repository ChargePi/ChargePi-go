package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export component settings of the ChargePi.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		exporter := cfg.GetExporter()
		config := viper.New()

		evseFlag := cmd.Flags().Lookup(settings.EvseFlag).Changed
		ocppFlag := cmd.Flags().Lookup(settings.OcppConfigPathFlag).Changed
		authFlag := cmd.Flags().Lookup(settings.AuthFileFlag).Changed

		// If the flag was set, export the EVSE configurations
		if evseFlag {
			err := cfg.ExportEVSEs(exporter, config, *evseFolderPath)
			if err != nil {
				return
			}
		}

		// If the flag was set, export the OCPP configuration
		if ocppFlag {
			err := cfg.ExportOcppConfiguration(exporter, config, *ocppConfigurationFilePath)
			if err != nil {
				return
			}
		}

		// If the flag was set, export tags.
		if authFlag {
			err := cfg.ExportLocalAuthList(exporter, config, *authFilePath)
			if err != nil {
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.
	exportCmd.PersistentFlags().StringVar(evseFolderPath, settings.EvseFlag, "./configs/evses", "evse folder path")
	exportCmd.PersistentFlags().StringVar(ocppConfigurationFilePath, settings.OcppConfigPathFlag, "./configs/settings.yaml", "OCPP config file path")
	exportCmd.PersistentFlags().StringVar(authFilePath, settings.AuthFileFlag, "./configs/authorization.yaml", "authorization file path")
}

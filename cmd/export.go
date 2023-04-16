package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	cfg "github.com/xBlaz3kx/ChargePi-go/internal/pkg/settings"
)

var (
	exportEvseFolderPath            *string
	exportOcppConfigurationFilePath *string
	exportAuthFilePath              *string
	exportSettingsFilePath          *string
)

// exportCommand represents the export command
func exportCommand() *cobra.Command {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export settings from ChargePi.",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			exporter := cfg.GetExporter()

			evseFlag := cmd.Flags().Lookup(settings.EvseFlag).Changed
			ocppFlag := cmd.Flags().Lookup(settings.OcppConfigPathFlag).Changed
			authFlag := cmd.Flags().Lookup(settings.AuthFileFlag).Changed
			settingsFlag := cmd.Flags().Lookup(settings.SettingsFlag).Changed

			// If the flag was set, export the EVSE configurations
			if evseFlag {
				err := exporter.ExportEVSESettingsToFile(*exportEvseFolderPath)
				if err != nil {
					return fmt.Errorf("could not export EVSE settings: %v", err)
				}
			}

			// If the flag was set, export the OCPP configuration
			if ocppFlag {
				err := exporter.ExportOcppConfigurationToFile(*exportOcppConfigurationFilePath)
				if err != nil {
					return fmt.Errorf("could not export OCPP configuration: %v", err)
				}
			}

			// If the flag was set, export tags.
			if authFlag {
				err := exporter.ExportLocalAuthListToFile(*exportAuthFilePath)
				if err != nil {
					return fmt.Errorf("could not export tags: %v", err)
				}
			}

			// If the flag was set, export settings.
			if settingsFlag {
				err := exporter.ExportChargePointSettingsToFile(*exportSettingsFilePath)
				if err != nil {
					return fmt.Errorf("could not export settings: %v", err)
				}
			}

			return nil
		},
	}

	exportEvseFolderPath = exportCmd.Flags().String(settings.EvseFlag, "./configs/evses", "evse folder path")
	exportOcppConfigurationFilePath = exportCmd.Flags().String(settings.OcppConfigPathFlag, "./configs/ocpp.yaml", "OCPP config file path")
	exportAuthFilePath = exportCmd.Flags().String(settings.AuthFileFlag, "./configs/authorization.yaml", "authorization file path")
	exportSettingsFilePath = exportCmd.Flags().String(settings.SettingsFlag, "./configs/settings.yaml", "settings file path")

	return exportCmd
}

package cmd

import (
	"fmt"

	"github.com/ChargePi/ChargePi-go/internal/auth"
	"github.com/ChargePi/ChargePi-go/internal/chargepoint"
	"github.com/ChargePi/ChargePi-go/internal/pkg/badger"
	cfg "github.com/ChargePi/ChargePi-go/internal/pkg/configuration"
	exporter2 "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/exporter"
	"github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/spf13/cobra"
)

var (
	exportEvseFolderPath            *string
	exportOcppConfigurationFilePath *string
	exportAuthFilePath              *string
	exportSettingsFilePath          *string
	databasePath                    *string
)

// exportCommand represents the export command
func exportCommand() *cobra.Command {
	exportCmd := &cobra.Command{
		Use:     "export",
		Short:   "Export settings from ChargePi.",
		Long:    ``,
		Version: chargepoint.FirmwareVersion,
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := badger.NewBadgerDb(*databasePath)
			if err != nil {
				return fmt.Errorf("could not create database: %v", err)
			}

			tagManager := auth.NewManager(db, db)

			configurationManager, err := ocpp_v16.NewV16ConfigurationManager(ocpp_v16.NewEmptyConfiguration())
			if err != nil {
				return fmt.Errorf("could not create OCPP configuration manager: %v", err)
			}

			settingsManager, err := manager.NewManager(db, db, configurationManager)
			if err != nil {
				return fmt.Errorf("could not create settings manager: %v", err)
			}

			exporter := exporter2.NewExporter(tagManager, settingsManager, db)

			evseFlag := cmd.Flags().Lookup(cfg.EvseFlag).Changed
			ocppFlag := cmd.Flags().Lookup(cfg.OcppConfigPathFlag).Changed
			authFlag := cmd.Flags().Lookup(cfg.AuthFileFlag).Changed
			settingsFlag := cmd.Flags().Lookup(cfg.SettingsFlag).Changed

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

	exportEvseFolderPath = exportCmd.Flags().String(cfg.EvseFlag, "./configs/evses", "evse folder path")
	exportOcppConfigurationFilePath = exportCmd.Flags().String(cfg.OcppConfigPathFlag, "./configs/ocpp.yaml", "OCPP config file path")
	exportAuthFilePath = exportCmd.Flags().String(cfg.AuthFileFlag, "./configs/authorization.yaml", "authorization file path")
	exportSettingsFilePath = exportCmd.Flags().String(cfg.SettingsFlag, "./configs/settings.yaml", "settings file path")
	databasePath = exportCmd.Flags().String("database", cfg.DatabasePath, "database path")

	return exportCmd
}

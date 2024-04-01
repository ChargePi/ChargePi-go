package cmd

import (
	chargePoint "github.com/ChargePi/ChargePi-go/internal/pkg/models/charge-point"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version of ChargePi",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("ChargePi version: %s", chargePoint.FirmwareVersion)
		},
	}
}

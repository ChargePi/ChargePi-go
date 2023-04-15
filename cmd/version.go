package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	chargePoint "github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/charge-point"
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

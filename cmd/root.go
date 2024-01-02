package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xBlaz3kx/ChargePi-go/internal/pkg/models/settings"
	"github.com/xBlaz3kx/ChargePi-go/pkg/observability/logging"
)

var rootCmd = &cobra.Command{
	Use:   "chargepi",
	Short: "ChargePi is an open-source Charge point project.",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(func() {
		logging.Setup(log.StandardLogger(), settings.Logging{}, viper.GetBool(settings.Debug))
	})

	rootCmd.AddCommand(runCommand())
	rootCmd.AddCommand(versionCommand())
	rootCmd.AddCommand(exportCommand())
	rootCmd.AddCommand(importCommand())

	rootCmd.PersistentFlags().BoolP(settings.DebugFlag, "d", false, "debug mode")
	_ = viper.BindPFlag(settings.Debug, rootCmd.PersistentFlags().Lookup(settings.DebugFlag))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.WithError(err).Fatal("Unable to run")
	}
}

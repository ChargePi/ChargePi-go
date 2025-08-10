package cmd

import (
	"github.com/ChargePi/ChargePi-go/internal/pkg/models/settings"
	"github.com/ChargePi/ChargePi-go/pkg/observability/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "chargepi",
	Short: "ChargePi is an open-source Charge point project.",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(func() {
		logger := logging.SetupZap(settings.Logging{}, viper.GetBool(settings.Debug))
		zap.ReplaceGlobals(logger)
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
		zap.L().Fatal("Unable to run", zap.Error(err))
	}
}

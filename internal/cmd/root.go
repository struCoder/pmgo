package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/struCoder/pmgo/internal/config"
	"github.com/struCoder/pmgo/internal/logger"
)

var (
	version   string
	buildTime string
	gitCommit string
	cfgFile   string
	cfg       *config.Config
	log       logger.Logger
)

// SetVersionInfo sets the version information for the application
func SetVersionInfo(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
}

// Execute runs the root command
func Execute(ctx context.Context, config *config.Config, logger logger.Logger) error {
	cfg = config
	log = logger

	rootCmd := &cobra.Command{
		Use:   "pmgo",
		Short: "A modern process manager for Go applications",
		Long: `PMGO is a lightweight, modern process manager written in Go for Go applications.
It helps you keep your applications alive forever, with features like auto-restart,
real-time monitoring, and web-based management interface.

Complete documentation is available at https://github.com/struCoder/pmgo`,
		Version: fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, gitCommit),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cfgFile != "" {
				return cfg.LoadFromFile(cfgFile)
			}
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.pmgo/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.Server.Host, "host", "localhost", "daemon host")
	rootCmd.PersistentFlags().IntVar(&cfg.Server.Port, "port", 9876, "daemon port")
	rootCmd.PersistentFlags().StringVar(&cfg.Logging.Level, "log-level", "info", "log level (debug, info, warn, error)")

	// Initialize viper
	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 9876)
	viper.SetDefault("log-level", "info")

	// Bind flags to viper
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	// Read config file if specified
	if configFile := viper.GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			log.Debugf("Config file not found or error reading: %v", err)
		}
	}

	// Add subcommands
	rootCmd.AddCommand(
		newServeCmd(),
		newStartCmd(),
		newStopCmd(),
		newRestartCmd(),
		newDeleteCmd(),
		newListCmd(),
		newInfoCmd(),
		newSaveCmd(),
		newKillCmd(),
		newWebCmd(),
		newLogsCmd(),
	)

	return rootCmd.ExecuteContext(ctx)
}

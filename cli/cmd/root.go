package cmd

import (
	"fmt"
	"os"

	"github.com/mrbrowser/mrbrowser/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	headless     bool
	chromiumPath string
	noSandbox    bool
	dbPath       string
	logLevel     string
)

// rootCmd is the base command.
var rootCmd = &cobra.Command{
	Use:   "mr-browser",
	Short: "Mr. Browser — intelligent browser automation engine",
	Long: `
  ███╗   ███╗██████╗     ██████╗ ██████╗  ██████╗ ██╗    ██╗███████╗███████╗██████╗
  ████╗ ████║██╔══██╗    ██╔══██╗██╔══██╗██╔═══██╗██║    ██║██╔════╝██╔════╝██╔══██╗
  ██╔████╔██║██████╔╝    ██████╔╝██████╔╝██║   ██║██║ █╗ ██║███████╗█████╗  ██████╔╝
  ██║╚██╔╝██║██╔══██╗    ██╔══██╗██╔══██╗██║   ██║██║███╗██║╚════██║██╔══╝  ██╔══██╗
  ██║ ╚═╝ ██║██║  ██║    ██████╔╝██║  ██║╚██████╔╝╚███╔███╔╝███████║███████╗██║  ██║
  ╚═╝     ╚═╝╚═╝  ╚═╝    ╚═════╝ ╚═╝  ╚═╝ ╚═════╝  ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝  ╚═╝

  Intent-driven, self-healing browser automation. No fragile selectors.`,
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.mrbrowser.yaml)")
	rootCmd.PersistentFlags().BoolVar(&headless, "headless", true, "run browser in headless mode")
	rootCmd.PersistentFlags().StringVar(&chromiumPath, "chromium", "", "path to Chromium executable")
	rootCmd.PersistentFlags().BoolVar(&noSandbox, "no-sandbox", false, "disable Chromium sandbox (required in Docker)")
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "./mrbrowser.db", "path to SQLite database")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level: debug, info, warn, error")

	_ = viper.BindPFlag("headless", rootCmd.PersistentFlags().Lookup("headless"))
	_ = viper.BindPFlag("chromium_path", rootCmd.PersistentFlags().Lookup("chromium"))
	_ = viper.BindPFlag("no_sandbox", rootCmd.PersistentFlags().Lookup("no-sandbox"))
	_ = viper.BindPFlag("db_path", rootCmd.PersistentFlags().Lookup("db"))
	_ = viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	// Register subcommands
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(screenshotCmd)
	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(debugCmd)
}

func initConfig() {
	viper.SetEnvPrefix("MRBROWSER")
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".mrbrowser")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}
	_ = viper.ReadInConfig()

	// Apply log level
	telemetry.SetLevelFromString(viper.GetString("log_level"))
}

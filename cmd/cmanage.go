package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// VERSION is set during build
	VERSION string
	// ConfigFile can be set for command
	ConfigFile string
)

// CmanageCmd is command's root command.
var CmanageCmd = &cobra.Command{
	Use:   "cmanage",
	Short: "cmanage manages Claroline Connect instances",
}

// Execute adds all child commands to the root command
func Execute(version string) {
	VERSION = version

	AddCommands()

	if err := CmanageCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
  CmanageCmd.PersistentFlags().StringVar(&ConfigFile, "config", "", "config file (default is /etc/cmanage/.cmanage.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	}

	viper.SetConfigName(".cmanage.toml") // name of config file (without extension)
	viper.AddConfigPath("/etc/cmanage/") // path to look for the config file in
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	/*if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}*/

	if err != nil {
		fmt.Println("No configuration file loaded - using defaults")
	}
}

// AddCommands adds child commands to the root command CmanageCmd.
func AddCommands() {
	CmanageCmd.AddCommand(versionCmd)
	CmanageCmd.AddCommand(releaseCmd)
	CmanageCmd.AddCommand(platformCmd)
}

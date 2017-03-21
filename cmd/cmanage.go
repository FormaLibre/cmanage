package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type config struct {
	domain             string
	releasePath        string
	webRoot            string
	mailjetPublicKey   string
	mailjetSecretKey   string
	mailjetFromName    string
	mailjetFromEmail   string
	mysqlHost          string
	mysqlPort          string
	mysqlRootUser      string
	mysqlRootPassword  string
	clarolinePassword  string
	clarolineUsername  string
	clarolineFirstName string
	clarolineLastName  string
	clarolineEmail     string
}

var (
	// VERSION is set during build
	VERSION string
	// ConfigFile can be set for command
	ConfigFile string
	// Config holder
	Config struct{config}
)

// RootCmd is command's root command.
var RootCmd = &cobra.Command{
	Use:   "cmanage",
	Short: "cmanage manages Claroline Connect instances",
}

// Execute adds all child commands to the root command
func Execute(version string) {
	VERSION = version

	addCommands()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	addCommands()
  RootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "", "config file (default is /etc/cmanage/cmanage.toml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigName("cmanage") // name of config file (without extension)
	viper.AddConfigPath("/etc/cmanage") // path to look for the config file in
	viper.ReadInConfig()

	if !viper.IsSet("main.domain") {
		fmt.Println("main.domain not set in config file")
		os.Exit(1)
	} else {
		Config.domain = viper.GetString("main.domain")
	}

	if !viper.IsSet("main.webroot") {
		fmt.Println("main.webroot not set in config file")
		os.Exit(1)
	} else {
		Config.webRoot = viper.GetString("main.webRoot")
	}

	if !viper.IsSet("main.releasePath") {
		fmt.Println("main.webroot not set in config file")
		os.Exit(1)
	} else {
		Config.releasePath = viper.GetString("main.releasePath")
	}

	if !viper.IsSet("mailjet.publicKey") {
		fmt.Println("mailjet.publicKey not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetPublicKey = viper.GetString("mailjet.publicKey")
	}

	if !viper.IsSet("mailjet.secretKey") {
		fmt.Println("mailjet.secretKey not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetSecretKey = viper.GetString("mailjet.secretKey")
	}

	if !viper.IsSet("mailjet.fromName") {
		fmt.Println("mailjet.fromName not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetFromName = viper.GetString("mailjet.fromName")
	}

	if !viper.IsSet("mailjet.fromEmail") {
		fmt.Println("mailjet.fromEmail not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetFromEmail = viper.GetString("mailjet.fromEmail")
	}

	if !viper.IsSet("mysql.host") {
		fmt.Println("mysql.host not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlHost = viper.GetString("mysql.host")
	}

	if !viper.IsSet("mysql.port") {
		fmt.Println("mysql.port not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlPort = viper.GetString("mysql.port")
	}

	if !viper.IsSet("mysql.rootUser") {
		fmt.Println("mysql.rootUser not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlRootUser = viper.GetString("mysql.rootUser")
	}

	if !viper.IsSet("mysql.rootPassword") {
		fmt.Println("mysql.rootPassword not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlRootPassword = viper.GetString("mysql.rootPassword")
	}

	if !viper.IsSet("claroline.password") {
		fmt.Println("claroline.password not set in config file")
		os.Exit(1)
	} else {
		Config.clarolinePassword = viper.GetString("claroline.password")
	}

	if !viper.IsSet("claroline.username") {
		fmt.Println("claroline.username not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineUsername = viper.GetString("claroline.username")
	}

	if !viper.IsSet("claroline.firstName") {
		fmt.Println("claroline.firstName not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineFirstName = viper.GetString("claroline.firstName")
	}

	if !viper.IsSet("claroline.lastName") {
		fmt.Println("claroline.lastName not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineLastName = viper.GetString("claroline.lastName")
	}

	if !viper.IsSet("claroline.email") {
		fmt.Println("claroline.email not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineEmail = viper.GetString("claroline.email")
	}
}

func addCommands() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(proxyCmd)
	RootCmd.AddCommand(releaseCmd)
	RootCmd.AddCommand(platformCmd)
}

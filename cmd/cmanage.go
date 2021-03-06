package cmd

import (
	"os"

	"github.com/ovh/go-ovh/ovh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/forma-libre/cmanage/utils"
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
	proxyAcmeEmail     string
	ovhAppKey          string
	ovhAppSecret       string
	ovhConsumerKey     string
	ovhRecordTarget    string
}

// Record defines the fields for an ovh domain record
type Record struct {
	Target    string `json:"target"`
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
}

var (
	// VERSION is set during build
	VERSION string
	// ConfigFile can be set for command
	ConfigFile string
	// Config holder
	Config struct{config}
	// OvhClient : OVH CLient
	OvhClient struct {ovh.Client}
)

// RootCmd is command's root command.
var RootCmd = &cobra.Command{
	Use:   "cmanage",
	Short: "cmanage manages Claroline Connect instances",
}

// Execute adds all child commands to the root command
func Execute(version string) {
	VERSION = version
	err = RootCmd.Execute();
	utils.Check(err)
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
		jww.ERROR.Println("main.domain not set in config file")
		os.Exit(1)
	} else {
		Config.domain = viper.GetString("main.domain")
	}

	if !viper.IsSet("main.webroot") {
		jww.ERROR.Println("main.webroot not set in config file")
		os.Exit(1)
	} else {
		Config.webRoot = viper.GetString("main.webRoot")
	}

	if !viper.IsSet("main.releasePath") {
		jww.ERROR.Println("main.webroot not set in config file")
		os.Exit(1)
	} else {
		Config.releasePath = viper.GetString("main.releasePath")
	}

	if !viper.IsSet("mailjet.publicKey") {
		jww.ERROR.Println("mailjet.publicKey not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetPublicKey = viper.GetString("mailjet.publicKey")
	}

	if !viper.IsSet("mailjet.secretKey") {
		jww.ERROR.Println("mailjet.secretKey not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetSecretKey = viper.GetString("mailjet.secretKey")
	}

	if !viper.IsSet("mailjet.fromName") {
		jww.ERROR.Println("mailjet.fromName not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetFromName = viper.GetString("mailjet.fromName")
	}

	if !viper.IsSet("mailjet.fromEmail") {
		jww.ERROR.Println("mailjet.fromEmail not set in config file")
		os.Exit(1)
	} else {
		Config.mailjetFromEmail = viper.GetString("mailjet.fromEmail")
	}

	if !viper.IsSet("mysql.host") {
		jww.ERROR.Println("mysql.host not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlHost = viper.GetString("mysql.host")
	}

	if !viper.IsSet("mysql.port") {
		jww.ERROR.Println("mysql.port not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlPort = viper.GetString("mysql.port")
	}

	if !viper.IsSet("mysql.rootUser") {
		jww.ERROR.Println("mysql.rootUser not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlRootUser = viper.GetString("mysql.rootUser")
	}

	if !viper.IsSet("mysql.rootPassword") {
		jww.ERROR.Println("mysql.rootPassword not set in config file")
		os.Exit(1)
	} else {
		Config.mysqlRootPassword = viper.GetString("mysql.rootPassword")
	}

	if !viper.IsSet("claroline.password") {
		jww.ERROR.Println("claroline.password not set in config file")
		os.Exit(1)
	} else {
		Config.clarolinePassword = viper.GetString("claroline.password")
	}

	if !viper.IsSet("claroline.username") {
		jww.ERROR.Println("claroline.username not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineUsername = viper.GetString("claroline.username")
	}

	if !viper.IsSet("claroline.firstName") {
		jww.ERROR.Println("claroline.firstName not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineFirstName = viper.GetString("claroline.firstName")
	}

	if !viper.IsSet("claroline.lastName") {
		jww.ERROR.Println("claroline.lastName not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineLastName = viper.GetString("claroline.lastName")
	}

	if !viper.IsSet("claroline.email") {
		jww.ERROR.Println("claroline.email not set in config file")
		os.Exit(1)
	} else {
		Config.clarolineEmail = viper.GetString("claroline.email")
	}

	if !viper.IsSet("proxy.acmeEmail") {
		jww.ERROR.Println("proxy.acmeEmail not set in config file")
		os.Exit(1)
	} else {
		Config.proxyAcmeEmail = viper.GetString("proxy.acmeEmail")
	}

	if !viper.IsSet("ovh.appKey") {
		jww.ERROR.Println("ovh.appKey not set in config file")
		os.Exit(1)
	} else {
		Config.ovhAppKey = viper.GetString("ovh.appKey")
	}

	if !viper.IsSet("ovh.appSecret") {
		jww.ERROR.Println("ovh.appSecret not set in config file")
		os.Exit(1)
	} else {
		Config.ovhAppSecret = viper.GetString("ovh.appSecret")
	}

	if !viper.IsSet("ovh.consumerKey") {
		jww.ERROR.Println("ovh.consumerKey not set in config file")
		os.Exit(1)
	} else {
		Config.ovhConsumerKey = viper.GetString("ovh.consumerKey")
	}

	if !viper.IsSet("ovh.recordTarget") {
		jww.ERROR.Println("ovh.recordTarget not set in config file")
		os.Exit(1)
	} else {
		Config.ovhRecordTarget = viper.GetString("ovh.recordTarget")
	}
}

func addCommands() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(proxyCmd)
	RootCmd.AddCommand(releaseCmd)
	RootCmd.AddCommand(platformCmd)
}

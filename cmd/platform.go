package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
	"strings"
	"database/sql"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/mailjet/mailjet-apiv3-go"

	"github.com/forma-libre/cmanage/bin"
	"github.com/forma-libre/cmanage/utils"
)

import _ "github.com/go-sql-driver/mysql" // This needs a comment, no idea why :)

var err error
var name string
var id string
var shellCommand string

// Main
var domain = viper.GetString("domain")
var webRoot = viper.GetString("webRoot")

// MySQL
var mysqlHost = viper.GetString("mysqlHost")
var mysqlPort = viper.GetString("mysqlPort")
var mysqlRootUser = viper.GetString("mysqlRootUser")
var mysqlRootPassword = viper.GetString("mysqlRootPassword")
var mysqlDsn = mysqlRootUser+":"+mysqlRootPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"

// Mailjet
var mailjetPublicKey = viper.GetString("mailjetPublicKey")
var mailjetSecretKey = viper.GetString("mailjetSecretKey")
var mailjetFromName = viper.GetString("mailjetFromName")
var mailjetFromEmail = viper.GetString("mailjetFromEmail")

// Claroline
var flPass = viper.GetString("flPass")
var flUsername = viper.GetString("flUsername")
var flFirstName = viper.GetString("flFirstName")
var flLastName = viper.GetString("flLastName")
var flEmail = viper.GetString("flEmail")

var platformCmd = &cobra.Command{
	Use:   "platform",
	Short: "Manage platforms",
}

var platformCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new platform",
	Run: func(cmd *cobra.Command, args []string) {

		subDomain := args[0]
		dbName := "cc_" + subDomain // This need filtering
		dbUser := dbName
		dbPass := utils.NewPassword(12)
		secret := utils.NewPassword(32)
		clientPass := utils.NewPassword(12)
		clientUsername := subDomain + "Admin"
		clientFirstName := "John" // Get this from flag
		clientLastName := "Doe" // Get this from flag
		clientEmail := "John.doe@client.net" // Get this from flag

		if len(args) == 0 {
			fmt.Println("You must specify the platform name.")
			fmt.Println("See 'platform create --help'")
			os.Exit(1)
		}
		if id == "" {
			fmt.Println("You must specify the platform ID.")
			fmt.Println("See 'platform create --help'")
			os.Exit(1)
		}
		if d, err := utils.Exists(webRoot + subDomain); d {
			fmt.Println("Error: " + subDomain + " directory allready exists.")
			utils.Check(err)
			os.Exit(1)
		}
		fmt.Println("Making directory " + subDomain)
		err = os.MkdirAll(webRoot+subDomain, 0755)
		utils.Check(err)

		fmt.Println("Making id file " + id)
		d := []byte(id)
		err = ioutil.WriteFile(webRoot+subDomain+"/id", d, 0644)
		utils.Check(err)

		fmt.Println("Creating docker-compose.yml")
		data, err := bin.Asset("data/docker-compose.yml")
		if err != nil {
			fmt.Println("Error: Asset not found")
		}
		err = ioutil.WriteFile(webRoot+subDomain+"/docker-compose.yml", data, 0644)
		utils.Check(err)

		fmt.Println("Copying Claroline (this can be long)")
		err = utils.CopyDir(releasePath+"current/claroline", webRoot+subDomain+"/claroline")
		utils.Check(err)

		fmt.Println("Generating .env file")
		e := []byte("SECRET=" + secret + "\nPLATFORM_SUBDOMAIN=" + subDomain + "\nPLATFORM_DOMAIN=" + domain + "\nDB_HOST=" + mysqlHost + "\nDB_USER=" + dbUser + "\nDB_PASSWORD=" + dbPass + "\nDB_NAME=" + dbName)
		err = ioutil.WriteFile(webRoot+subDomain+"/.env", e, 0644)
		utils.Check(err)

		db, err := sql.Open("mysql", mysqlDsn)
		utils.Check(err)
		defer db.Close()

		stm := "CREATE USER '" + dbUser + "'@'localhost' IDENTIFIED BY '" + dbPass + "'"
		_, err = db.Exec(stm)
		utils.Check(err)

		stm = "CREATE USER '" + dbUser + "'@'%' IDENTIFIED BY '" + dbPass + "'"
		_, err = db.Exec(stm)
		utils.Check(err)

		stm = "CREATE DATABASE " + dbName
		_, err = db.Exec(stm)
		utils.Check(err)

		stm = "GRANT ALL ON " + dbName + ".* TO '" + dbUser + "'@'localhost'"
		_, err = db.Exec(stm)
		utils.Check(err)

		stm = "GRANT ALL ON " + dbName + ".* TO '" + dbUser + "'@'%'"
		_, err = db.Exec(stm)
		utils.Check(err)

		fmt.Println("Creating deploy script")
		data, err = bin.Asset("data/deploy.sh")
		if err != nil {
			fmt.Println("Error: Asset not found")
		}

		err = os.Chmod(webRoot+subDomain+"/claroline/app/cache", 0777)
		utils.Check(err)

		err = os.Chmod(webRoot+subDomain+"/claroline/app/config", 0777)
		utils.Check(err)

		err = os.Chmod(webRoot+subDomain+"/claroline/app/logs", 0777)
		utils.Check(err)

		err = os.Chmod(webRoot+subDomain+"/claroline/app/sessions", 0777)
		utils.Check(err)

		err = os.Chmod(webRoot+subDomain+"/claroline/web/uploads", 0777)
		utils.Check(err)

		err = os.Chmod(webRoot+subDomain+"/claroline/files", 0777)
		utils.Check(err)

		shellCommand = "docker pull claroline/claroline-docker:prod"
		out, err := exec.Command("sh", "-c", shellCommand).Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

    env := "SECRET=\""+secret+"\" PLATFORM_SUBDOMAIN=\""+subDomain+"\" PLATFORM_DOMAIN=\""+domain+"\" DB_HOST=\""+mysqlHost+"\" DB_USER=\""+dbUser+"\" DB_PASSWORD=\""+dbPass+"\" DB_NAME=\""+dbName+"\""

		shellCommand = "cd "+webRoot+subDomain+"/ && env "+env+" docker stack deploy --compose-file docker-compose.yml "+subDomain
		fmt.Println("Executing: "+shellCommand)
		out, err = exec.Command("sh", "-c", shellCommand).Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

		fmt.Println("Sleeping 20s") // TODO this coud be better
		time.Sleep(20000 * time.Millisecond)

		out, err = exec.Command("docker", "ps", "-q", "--filter", "name="+subDomain+"_claroline").Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

		containerID := strings.TrimSpace(string(out))

		cmdStr :=  "docker exec -i " + containerID + " sh -c 'cd claroline && php scripts/configure.php'"
		out, err = exec.Command("sh", "-c", cmdStr).Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

		fmt.Println("Composer install")
		cmdStr = "docker exec -i " + containerID + " sh -c 'cd claroline && composer install'"
		out, err = exec.Command("sh", "-c", cmdStr).Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

		fmt.Println("Composer fast-install")
		cmdStr = "docker exec -i " + containerID + " sh -c 'cd claroline && composer fast-install'"
		out, err = exec.Command("sh", "-c", cmdStr).Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

		fmt.Println("Create Forma-Libre User")
		cmdStr = "docker exec -i " + containerID + " sh -c 'cd claroline && php app/console claroline:user:create -a " + flFirstName + " " + flLastName + " " + flUsername + " " + flPass + " " + flEmail + "'"
		out, err = exec.Command("sh", "-c", cmdStr).Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

		fmt.Println("Create Client User")
		cmdStr = "docker exec -i " + containerID + " sh -c 'cd claroline && php app/console claroline:user:create -a " + clientFirstName + " " + clientLastName + " " + clientUsername + " " + clientPass + " " + clientEmail + "'"
		fmt.Println(cmdStr)
		out, err = exec.Command("sh", "-c", cmdStr).Output()
		utils.Check(err)
		fmt.Printf("%s\n", out)

		mj := mailjet.NewMailjetClient(mailjetPublicKey, mailjetSecretKey)

    param := &mailjet.InfoSendMail{
        FromEmail: mailjetFromEmail,
        FromName: mailjetFromName,
        Recipients: []mailjet.Recipient{
            mailjet.Recipient{
                Email: mailjetFromEmail,
            },
        },
        Subject: "Your new Claroline Connect platform has been created!",
        TextPart: "Hello [[USERNAME ]]\nYour new Claroline Connect platform is ready, here is the information you need to connect:\nURL: https://"+ subDomain + "." + domain + "\nUsername: " + clientUsername + "\nPassword: " + clientPass +"\n" + "Enjoy!", // TODO this could be better
    }

    res, err := mj.SendMail(param)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(res)
    }

		fmt.Println(id)
	},
}

var platformLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List platforms",
	Run: func(cmd *cobra.Command, args []string) {
		var id string
		files, _ := ioutil.ReadDir(webRoot)
		table := uitable.New()
		table.MaxColWidth = 80
		table.AddRow("PLATFORM ID", "NAME", "CREATED")
		for _, f := range files {
			if f.IsDir() {
				if _, err := os.Stat(webRoot + f.Name() + "/id"); err == nil {
					b, err := ioutil.ReadFile(webRoot + f.Name() + "/id")
					utils.Check(err)
					id = string(b)
				} else {
					id = "Undefined"
				}
				table.AddRow(id, f.Name(), humanize.Time(f.ModTime()))
			}
		}
		fmt.Println(table)
	},
}

var platformRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove one or more platforms",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify the platform name.")
			fmt.Println("See 'platform rm --help'")
			os.Exit(1)
		}
		for _, v := range args {
			os.RemoveAll(webRoot + v)
			subDomain := v
			dbName := "cc_" + subDomain //This need filtering

			db, err := sql.Open("mysql", mysqlDsn)
			utils.Check(err)
			defer db.Close()

			fmt.Println("Removing database")
			var stm string
			dbUser := dbName

			stm = "DROP DATABASE IF EXISTS "+dbName
			_, err = db.Exec(stm)
	    utils.Check(err)

			fmt.Println("Removing database users")
			stm = "DROP USER IF EXISTS " + dbUser + "@localhost"
			_, err = db.Exec(stm)
	    utils.Check(err)

			stm = "DROP USER IF EXISTS " + dbUser
			_, err = db.Exec(stm)
	    utils.Check(err)

			fmt.Println(v)
		}

	},
}

var platformStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start one or more platforms",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Start platform command")
	},
}

var platformStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop one or more platforms",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Stop platform command")
	},
}

func init() {
	platformCmd.AddCommand(platformCreateCmd)
	platformCmd.AddCommand(platformLsCmd)
	platformCmd.AddCommand(platformRmCmd)
	platformCmd.AddCommand(platformStartCmd)
	platformCmd.AddCommand(platformStopCmd)
	platformCreateCmd.Flags().StringVarP(&name, "name", "n", "", "Claroline Connect Platform Name")
	platformCreateCmd.Flags().StringVarP(&id, "id", "i", "", "Forma Libre Manager ID")
}

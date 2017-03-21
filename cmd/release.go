package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/gosuri/uitable"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"

	"github.com/forma-libre/cmanage/utils"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Manage releases",
}

var releaseLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List releases",
	Run: func(cmd *cobra.Command, args []string) {
		files, _ := ioutil.ReadDir(Config.releasePath)
		table := uitable.New()
		table.MaxColWidth = 80
		table.AddRow("NAME", "CREATED", "SIZE")
		for _, f := range files {
			if f.IsDir() {
				dir := Config.releasePath + f.Name()
				info, err := os.Lstat(dir)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				table.AddRow(f.Name(), humanize.Time(f.ModTime()), humanize.Bytes(uint64(utils.DiskUsage(dir, info))))
			}
		}
		fmt.Println(table)
	},
}

var releaseGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get release",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify the release.")
			fmt.Println("See 'release get --help'")
			os.Exit(1)
		}
		if _, err := os.Stat(Config.releasePath); os.IsNotExist(err) {
			os.MkdirAll(Config.releasePath, 0755)
		}
		utils.Check(err)
		resp, err := http.Get("http://packages.claroline.net/releases/" + args[0] + "/claroline-" + args[0] + ".tar.gz")
		if err == nil && (resp.StatusCode == 200) {
			fmt.Println("Downloading release " + args[0] + " archive.")
			defer resp.Body.Close()
			out, err := os.Create(Config.releasePath + "/claroline-" + args[0] + ".tar.gz")
			defer out.Close()
			utils.Check(err)
			_, err = io.Copy(out, resp.Body)
			fmt.Println("Unpacking release " + args[0] + ".")
			err = archiver.TarGz.Open(Config.releasePath+"/claroline-"+args[0]+".tar.gz", Config.releasePath+args[0])
			utils.Check(err)
			fmt.Println("Removing release " + args[0] + " archive.")
			err = os.Remove(Config.releasePath + "/claroline-" + args[0] + ".tar.gz")
			utils.Check(err)
		} else {
			utils.Check(err)
			if resp.StatusCode == 404 {
				fmt.Println("The release could not be found.")
			}
			os.Exit(1)
		}
	},
}

var releaseSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the current releases",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify the release.")
			fmt.Println("See 'release set --help'")
			os.Exit(1)
		}
		path := Config.releasePath + args[0]
		if d, err := utils.NotExists(path); d {
			utils.Check(err)
			fmt.Println("Error : Non existant release.")
			fmt.Println("See 'release ls'")
			fmt.Println("See 'release set --help'")
			os.Exit(1)
		}
		os.Symlink(path, Config.releasePath+"current")
		fmt.Println("Current release set command to " + args[0] + ".")
	},
}

var releaseRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove one or more releases",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must specify the release.")
			fmt.Println("See 'release rm --help'")
			os.Exit(1)
		}
		for _, v := range args {
			os.RemoveAll(Config.releasePath + v)
			fmt.Println(v)
		}
	},
}

func init() {
	releaseCmd.AddCommand(releaseLsCmd)
	releaseCmd.AddCommand(releaseRmCmd)
	releaseCmd.AddCommand(releaseGetCmd)
	releaseCmd.AddCommand(releaseSetCmd)
}

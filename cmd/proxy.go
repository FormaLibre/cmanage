package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/forma-libre/cmanage/utils"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Manage proxy",
}

var proxyCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create cmanage proxy",
	Run: func(cmd *cobra.Command, args []string) {
		shellCommand = "docker network create --driver=overlay --attachable cmanager_proxy_network"
		out, _ := exec.Command("sh", "-c", shellCommand).Output()
		utils.Check(err)
		fmt.Printf("%s", out)

		shellCommand = "docker service create --name cmanager_proxy --publish 80:80 --publish 8080:8080 --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock --network cmanager_proxy_network traefik -c /dev/null --docker --docker.swarmmode --docker.domain="+Config.domain+" --docker.watch --web"
		out, _ = exec.Command("sh", "-c", shellCommand).Output()
		utils.Check(err)
		fmt.Printf("%s", out)
	},
}

var proxyRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove cmanage proxy",
	Run: func(cmd *cobra.Command, args []string) {
		shellCommand = "docker network rm cmanager_proxy_network"
		out, _ := exec.Command("sh", "-c", shellCommand).Output()
		utils.Check(err)
		fmt.Printf("%s", out)

		shellCommand = "docker service rm cmanager_proxy"
		out, _ = exec.Command("sh", "-c", shellCommand).Output()
		utils.Check(err)
		fmt.Printf("%s", out)
	},
}

func init() {
	proxyCmd.AddCommand(proxyCreateCmd)
	proxyCmd.AddCommand(proxyRmCmd)
}

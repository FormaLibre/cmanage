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
		shellCommand = "docker service create --name cmanager_proxy --publish 443:443 --publish 80:80 --publish 8080:8080 --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock --mount type=bind,source=/var/tmp,target=/etc/traefik/acme --network cmanager_proxy_network traefik --entryPoints='Name:http Address::80 Redirect.EntryPoint:https' --entryPoints='Name:https Address::443 TLS' --defaultEntryPoints=http,https --acme.entryPoint=https --acme.email="+Config.proxyAcmeEmail+" --acme.storage=/etc/traefik/acme/acme.json --acme.domains="+Config.domain+" --acme.onHostRule=true  --docker --docker.swarmmode --docker.domain="+Config.domain+" --docker.watch --web"
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

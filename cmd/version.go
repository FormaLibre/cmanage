package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show cmanage client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cleint:")
		fmt.Println(" Version: " + VERSION)
	},
}

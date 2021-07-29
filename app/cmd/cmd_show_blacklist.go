package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showBlacklistCmd = &cobra.Command{
	Use:   "show-blacklist",
	Short: "Show blacklist",
	Long:  `Show service blacklist`,
	Run:   showBlacklistExecute,
}

func init() {
	RootCmd.AddCommand(showBlacklistCmd)
}

func showBlacklistExecute(cmd *cobra.Command, args []string) {
	fmt.Println("blacklist show called")
}

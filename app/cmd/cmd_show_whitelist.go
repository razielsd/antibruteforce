package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var showWhitelistCmd = &cobra.Command{
	Use:   "show-whitelist",
	Short: "Show whitelist",
	Long:  `Show service whitelist`,
	Run:   showWhitelistExecute,
}

func init() {
	RootCmd.AddCommand(showWhitelistCmd)
}

func showWhitelistExecute(cmd *cobra.Command, args []string) {
	fmt.Println("whitelist show called")
}

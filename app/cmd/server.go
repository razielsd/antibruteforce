package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run service",
	Long:  `Run service`,
	Run:   serverExecute,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func serverExecute(cmd *cobra.Command, args []string) {
	fmt.Println("server called")
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/razielsd/antibruteforce/app/api"
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

func serverExecute(command *cobra.Command, args []string) {
	abf, err := api.NewAbfAPI(abfConfig, abfLogger)
	if err != nil {
		fmt.Printf("Error starting service: %s\n", err)
		os.Exit(1)
	}
	abf.Run()
}

package cmd

import (
	"github.com/razielsd/antibruteforce/app/api"
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

func serverExecute(command *cobra.Command, args []string) {
	abf := api.NewAbfAPI(abfConfig, abfLogger)
	abf.Run()
}

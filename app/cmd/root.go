package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	abfconfig "github.com/razielsd/antibruteforce/app/config"
	"github.com/razielsd/antibruteforce/app/logger"
)

var RootCmd = &cobra.Command{
	Use:   "antibruteforce",
	Short: "antibruteforce",
	Long:  `Antibruteforce service cli`,
}

var (
	abfLogger *zap.Logger
	abfConfig abfconfig.AppConfig
)

func Execute() {
	initEnv()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initEnv() {
	var err error
	abfConfig, err = abfconfig.GetConfig()
	if err != nil {
		fmt.Printf("Unable get config: %s\n", err)
		os.Exit(1)
	}
	abfLogger, err = logger.GetLogger(abfConfig)
	if err != nil {
		fmt.Printf("Unable init logger: %s\n", err)
		os.Exit(1)
	}
}

func extractFirstArg(cmd *cobra.Command, args []string, errMessage string) string {
	if len(args) < 1 {
		fmt.Println(errMessage)
		_ = cmd.Usage()
		os.Exit(1)
	}
	return args[0]
}

func printCli(data string, err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println(data)
	}
}

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	abfconfig "github.com/razielsd/antibruteforce/app/config"
	"github.com/razielsd/antibruteforce/app/logger"
)

var rootCmd = &cobra.Command{
	Use:   "antibruteforce",
	Short: "antibruteforce",
	Long:  `Antibruteforce service cli`,
}

// Execute main run point.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getConfigOrDie() abfconfig.AppConfig {
	abfConfig, err := abfconfig.GetConfig()
	if err != nil {
		fmt.Printf("Unable load config: %s\n", err)
		os.Exit(1)
	}
	return abfConfig
}

func getLoggerOrDie(cfg abfconfig.AppConfig) *zap.Logger {
	abfLogger, err := logger.GetLogger(cfg)
	if err != nil {
		fmt.Printf("Unable init logger: %s\n", err)
		os.Exit(1)
	}
	return abfLogger
}

func extractFirstArg(args []string) (string, error) {
	if len(args) < 1 {
		return "", errors.New("empty args")
	}
	return args[0], nil
}

func extractFirstArgOrDie(cmd *cobra.Command, args []string, errMessage string) string {
	param, err := extractFirstArg(args)
	if err != nil {
		fmt.Println(errMessage)
		_ = cmd.Usage()
		os.Exit(1)
	}
	return param
}

func printCli(data string, err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println(data)
	}
}

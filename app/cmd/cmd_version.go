package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	buildDate = "0000-00-00"
	gitHash   = "undefined"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Long:  `Show version`,
	Run:   versionExecute,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionExecute(cmd *cobra.Command, args []string) {
	fmt.Println("Version:\t", version)
	fmt.Println("Build date:\t", buildDate)
	fmt.Println("Git hash:\t", gitHash)
}

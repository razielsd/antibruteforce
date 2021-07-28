package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var dropBucketCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop bucket for login, password or ip/mask",
	Long:  `Drop bucket from service`,
	Run: dropBucketExecute,
}

func init() {
	RootCmd.AddCommand(dropBucketCmd)
}

func dropBucketExecute(cmd *cobra.Command, args []string) {
	fmt.Println("drop bucket called")
}
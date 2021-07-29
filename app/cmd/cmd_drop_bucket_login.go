package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dropBucketLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Drop bucket for login",
	Long:  `Drop bucket for login`,
	Run:   dropBucketLoginExecute,
}

func dropBucketLoginExecute(cmd *cobra.Command, args []string) {
	fmt.Println("drop bucket by login called")
}

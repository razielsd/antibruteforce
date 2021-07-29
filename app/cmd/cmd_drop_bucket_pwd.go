package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dropBucketPwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Drop bucket for password",
	Long:  `Drop bucket for password`,
	Run:   dropBucketPwdExecute,
}

func dropBucketPwdExecute(cmd *cobra.Command, args []string) {
	fmt.Println("drop bucket by password called")
}

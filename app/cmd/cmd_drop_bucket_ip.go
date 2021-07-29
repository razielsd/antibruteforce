package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var dropBucketIPCmd = &cobra.Command{
	Use:   "ip",
	Short: "Drop bucket for ip/mask",
	Long:  `Drop bucket for ip/mask`,
	Run:   dropBucketIPExecute,
}

func dropBucketIPExecute(cmd *cobra.Command, args []string) {
	fmt.Println("drop bucket by ip called")
}

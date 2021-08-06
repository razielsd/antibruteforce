package cmd

import (
	"github.com/spf13/cobra"

	"github.com/razielsd/antibruteforce/app/cli"
)

var rootBucketCmd = &cobra.Command{
	Use:   "bucket",
	Short: "Drop bucket by login, password or ip",
	Long:  "Drop bucket by login, password or ip",
}

var dropBucketCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop bucket by login, password or ip",
	Long:  "Drop bucket by login, password or ip",
}

var dropBucketLoginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Drop bucket by login",
	Long:    `Drop bucket by login`,
	Run:     dropBucketByLogin,
	Example: "abf bucket drop login <login>",
}

var dropBucketPwdCmd = &cobra.Command{
	Use:     "pwd",
	Short:   "Drop bucket by password",
	Long:    `Drop bucket by password`,
	Run:     dropBucketByPwd,
	Example: "abf bucket drop pwd <password>",
}

var dropBucketIPCmd = &cobra.Command{
	Use:     "ip",
	Short:   "Drop bucket by IP",
	Long:    `Drop bucket by IP`,
	Run:     dropBucketByIP,
	Example: "abf bucket drop ip <ip>",
}

func init() {
	dropBucketCmd.AddCommand(dropBucketLoginCmd)
	dropBucketCmd.AddCommand(dropBucketPwdCmd)
	dropBucketCmd.AddCommand(dropBucketIPCmd)
	rootBucketCmd.AddCommand(dropBucketCmd)
	RootCmd.AddCommand(rootBucketCmd)
}

func dropBucketByLogin(cmd *cobra.Command, args []string) {
	key := extractFirstArg(cmd, args, "Require login")
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.DropBucketByLogin(key)
	printCli(message, err)
}

func dropBucketByPwd(cmd *cobra.Command, args []string) {
	key := extractFirstArg(cmd, args, "Require password")
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.DropBucketByPwd(key)
	printCli(message, err)
}

func dropBucketByIP(cmd *cobra.Command, args []string) {
	key := extractFirstArg(cmd, args, "Require IP")
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.DropBucketByIP(key)
	printCli(message, err)
}

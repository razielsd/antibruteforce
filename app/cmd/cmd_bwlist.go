package cmd

import (
	"github.com/spf13/cobra"

	"github.com/razielsd/antibruteforce/app/cli"
)

var rootBlacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "Show blacklist",
	Long:  `Show service blacklist`,
}

var showBlacklistCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show blacklist",
	Long:    `Show service blacklist`,
	Run:     showBlacklistExecute,
	Example: "abf blacklist show",
}

var appendBlacklistCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add ip/mask to blacklist",
	Long:    `Add ip/mask to blacklist`,
	Run:     addBlacklistExecute,
	Example: "abf blacklist add <ip/mask>",
}

var rmBlacklistCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove ip/mask from blacklist",
	Long:    `Remove ip/mask service blacklist`,
	Run:     rmBlacklistExecute,
	Example: "abf blacklist rm <ip/mask>",
}

var rootWhitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Show whitelist",
	Long:  `Show service whitelist`,
}

var showWhitelistCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show whitelist",
	Long:    `Show service whitelist`,
	Run:     showWhitelistExecute,
	Example: "abf whitelist show",
}

var appendWhitelistCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add ip/mask to whitelist",
	Long:    `Add ip/mask to whitelist`,
	Run:     addWhitelistExecute,
	Example: "abf whitelist add <ip/mask>",
}

var rmWhitelistCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove ip/mask from whitelist",
	Long:    `Remove ip/mask service whitelist`,
	Run:     rmWhitelistExecute,
	Example: "abf whitelist rm <ip/mask>",
}

func init() {
	rootBlacklistCmd.AddCommand(showBlacklistCmd)
	rootBlacklistCmd.AddCommand(appendBlacklistCmd)
	rootBlacklistCmd.AddCommand(rmBlacklistCmd)
	RootCmd.AddCommand(rootBlacklistCmd)

	rootWhitelistCmd.AddCommand(showWhitelistCmd)
	rootWhitelistCmd.AddCommand(appendWhitelistCmd)
	rootWhitelistCmd.AddCommand(rmWhitelistCmd)
	RootCmd.AddCommand(rootWhitelistCmd)
}

func showBlacklistExecute(cmd *cobra.Command, args []string) {
	cliClient := cli.NewCli()
	message, err := cliClient.ShowBlacklist()
	printCli(message, err)
}

func addBlacklistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArg(cmd, args, "Require ip/mask")
	cliClient := cli.NewCli()
	message, err := cliClient.AppendBlacklist(ip)
	printCli(message, err)
}

func rmBlacklistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArg(cmd, args, "Require ip/mask")
	cliClient := cli.NewCli()
	message, err := cliClient.RemoveBlacklist(ip)
	printCli(message, err)
}

func showWhitelistExecute(cmd *cobra.Command, args []string) {
	cliClient := cli.NewCli()
	message, err := cliClient.ShowWhitelist()
	printCli(message, err)
}

func addWhitelistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArg(cmd, args, "Require ip/mask")
	cliClient := cli.NewCli()
	message, err := cliClient.AppendWhitelist(ip)
	printCli(message, err)
}

func rmWhitelistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArg(cmd, args, "Require ip/mask")
	cliClient := cli.NewCli()
	message, err := cliClient.RemoveWhitelist(ip)
	printCli(message, err)
}

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/razielsd/antibruteforce/app/cli"
)

var rootBlacklistCmd = &cobra.Command{
	Use:   "blacklist",
	Short: "Show/add/remove blacklist",
	Long:  `Show/add/remove blacklist`,
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
	Short:   "Add ip/subnet to blacklist",
	Long:    `Add ip/subnet to blacklist`,
	Run:     addBlacklistExecute,
	Example: "abf blacklist add <ip/subnet>",
}

var rmBlacklistCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove ip/subnet from blacklist",
	Long:    `Remove ip/subnet service blacklist`,
	Run:     rmBlacklistExecute,
	Example: "abf blacklist rm <ip/subnet>",
}

var rootWhitelistCmd = &cobra.Command{
	Use:   "whitelist",
	Short: "Show/add/remove whitelist",
	Long:  `Show/add/remove whitelist`,
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
	Short:   "Add ip/subnet to whitelist",
	Long:    `Add ip/subnet to whitelist`,
	Run:     addWhitelistExecute,
	Example: "abf whitelist add <ip/subnet>",
}

var rmWhitelistCmd = &cobra.Command{
	Use:     "rm",
	Short:   "Remove ip/subnet from whitelist",
	Long:    `Remove ip/subnet service whitelist`,
	Run:     rmWhitelistExecute,
	Example: "abf whitelist rm <ip/subnet>",
}

func init() {
	rootBlacklistCmd.AddCommand(showBlacklistCmd)
	rootBlacklistCmd.AddCommand(appendBlacklistCmd)
	rootBlacklistCmd.AddCommand(rmBlacklistCmd)
	rootCmd.AddCommand(rootBlacklistCmd)

	rootWhitelistCmd.AddCommand(showWhitelistCmd)
	rootWhitelistCmd.AddCommand(appendWhitelistCmd)
	rootWhitelistCmd.AddCommand(rmWhitelistCmd)
	rootCmd.AddCommand(rootWhitelistCmd)
}

func showBlacklistExecute(cmd *cobra.Command, args []string) {
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.ShowBlacklist()
	printCli(message, err)
}

func addBlacklistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArgOrDie(cmd, args, "Require ip/subnet")
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.AppendBlacklist(ip)
	printCli(message, err)
}

func rmBlacklistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArgOrDie(cmd, args, "Require ip/subnet")
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.RemoveBlacklist(ip)
	printCli(message, err)
}

func showWhitelistExecute(cmd *cobra.Command, args []string) {
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.ShowWhitelist()
	printCli(message, err)
}

func addWhitelistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArgOrDie(cmd, args, "Require ip/subnet")
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.AppendWhitelist(ip)
	printCli(message, err)
}

func rmWhitelistExecute(cmd *cobra.Command, args []string) {
	ip := extractFirstArgOrDie(cmd, args, "Require ip/subnet")
	cliClient := cli.NewCli(getConfigOrDie())
	message, err := cliClient.RemoveWhitelist(ip)
	printCli(message, err)
}

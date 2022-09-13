package main

import (
	"github.com/spf13/cobra"
	"os"
)

const ISO8601 = "2006-01-02"

var listAliases = []string{"ls", "l"}

func main() {
	habitCmd.AddCommand(
		habitListCmd,
		habitsAddCmd,
	)

	rootCmd.AddCommand(
		authCmd,
		habitCmd,
		logCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Global flags:
//   - environment
var rootCmd = &cobra.Command{
	Use:   "ape <command> <subcommand> [flags]",
	Short: "Show help when using this app",
	Long:  "Curious Ape command line client.",
}

func init() {
	// Define flags and configuration settings
	rootCmd.PersistentFlags().StringP("host", "e", "https://ape.danicos.me/api", "Server host:port combination")
	rootCmd.PersistentFlags().String("config-dir-path", "/home/daniel/.ape/cli", "Configuration files location")

	// habits
	habitListCmd.Flags().StringP("by", "b", "month", "by which unit to list: day/week/month etc")
	habitListCmd.Flags().StringP("from", "f", "current", "from which period: current, previous, next")
	// auth
	authCmd.Flags().StringP("username", "u", "", "Username")
	authCmd.Flags().StringP("password", "p", "", "Password")
	authCmd.MarkFlagRequired("username")
	authCmd.MarkFlagRequired("password")
}

// Cobra

// List today
// List week
// List month

// Add Habit
// 		All
// Add Record
// 		Fitness

// What commands to I want
// ape habits ls --by=month/week/day

// ape logs sleep/fitness/work ls --by=month/week/day

// ape days ls --by=month/week/day

// ape ls habits/days/

// root command will default to listing the help.

// route /home/$USER/.ape/cli/rc.json
// ape login
//      if it is not logged-in do nothing.
// 		when logging in try to make a request 	if 200 you are logged in.
//												if 400+ you are not logged in.
// 		username ->
// 		password ->

// ape day/d 	ls
// ape log/l 	ls add/a (only fitness for now).

// ls -by=day/week/month/year/all (defaults to current day)

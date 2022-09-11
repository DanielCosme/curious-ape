package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:     "log [resource to list]",
	Short:   "Manage activity logs like sleep, fitness, etc",
	Long:    "List, create, update and delete logs",
	Aliases: []string{"logs", "ls", "l"},
	Run: func(cmd *cobra.Command, args []string) {
		// ape habit/h 	ls add/a

		// Verify login -> if it fails return error prompting the user to login.
		// ls/s -> subcommand

		// add/a -> subcommand
		fmt.Println("not implemented")
	},
}

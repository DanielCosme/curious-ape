package main

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/client"
	"github.com/spf13/cobra"
	"strings"
)

var habitCmd = &cobra.Command{
	Use:               "habit <command> [flags]",
	Short:             "Manage habits",
	Long:              "List, create, update and delete habit logs",
	Aliases:           []string{"habits", "h"},
	ArgAliases:        []string{},
	PersistentPreRunE: loadCredentials,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("I will default to list the habits for today")
	},
}

var habitListCmd = &cobra.Command{
	Use:     "list",
	Short:   "habit list",
	Long:    `habit list`,
	Aliases: listAliases,
	// Aliases:    []string{"habits", "h"},
	ArgAliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing Habits")

		err := client.Ping()
		CheckErr(err)

		by, _ := cmd.Flags().GetString("by")
		from, _ := cmd.Flags().GetString("from")
		fmt.Println("Listing", strings.Join(args, ", "))
		fmt.Println("By", by)
		fmt.Println("From", from)
	},
}

package main

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/client"
	"github.com/spf13/cobra"
	"time"
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
	Use:        "list",
	Short:      "habit list",
	Long:       `habit list`,
	Aliases:    listAliases,
	ArgAliases: []string{},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing Habits")

		c := client.DefaultClient
		hs, err := c.Habits.List()
		CheckErr(err)

		for _, h := range hs {
			fmt.Printf("Date: %s - Status: %s - Type: %s\n", h.Date.Format(time.Stamp), h.Status, h.Type.Str())
			for _, l := range h.Logs {
				fmt.Printf(" ---> Success: %v - Origin: %s\n", l.Success, l.Origin)
			}
		}

		// Get all habits:
		// 		By:
		// 		From:

		// by, _ := cmd.Flags().GetString("by")
		// from, _ := cmd.Flags().GetString("from")
		// fmt.Println("Listing", strings.Join(args, ", "))
		// fmt.Println("By", by)
		// fmt.Println("From", from)
	},
}

func period(from time.Time, period, by string) time.Time {
	return time.Now()
}

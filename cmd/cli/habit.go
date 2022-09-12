package main

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/client"
	"github.com/danielcosme/go-sdk/dates"
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
		by, _ := cmd.Flags().GetString("by")     // endDate -> 	day, week, month, year
		from, _ := cmd.Flags().GetString("from") // startDate -> current, previous

		hs, err := client.DefaultClient.Habits.List(period(by, from))
		CheckErr(err)
		if len(hs) == 0 {
			fmt.Println("No habits found.")
			return
		}
		mapHabitsByDay := map[time.Time][]types.HabitTransport{}

		for _, h := range hs {
			mapHabitsByDay[*h.Date] = append(mapHabitsByDay[*h.Date], h)
		}
		for d, hs := range mapHabitsByDay {
			fmt.Printf("Day: %s\n", d.Format(time.RubyDate))
			for _, h := range hs {
				fmt.Printf(" --> %s: %s\n", h.Type.Str(), h.Status)
				// for _, l := range h.Logs {
				// 	fmt.Printf("     Success: %v - Origin: %s\n", l.Success, l.Origin)
				// }
			}
			fmt.Println()
		}
	},
}

func period(by, from string) (time.Time, time.Time) {
	var startDate, endDate time.Time
	switch by {
	case "day":
		switch from {
		case "current":
			startDate = time.Now().Local()
			endDate = dates.EndOfDay(startDate)
		case "previous":
			startDate = time.Now().AddDate(0, 0, -1)
			endDate = dates.EndOfDay(startDate)
		}
	case "week":
		switch from {
		case "current":
			startDate = dates.BeginningOfWeek(time.Now())
			endDate = dates.EndOfWeek(time.Now())
		case "previous":
			startDate = dates.BeginningOfWeek(time.Now()).AddDate(0, 0, -7)
			endDate = dates.EndOfWeek(time.Now()).AddDate(0, 0, -7)
		}
	case "month":
	case "year":
	default:
		// All
	}
	fmt.Println("Start: ", startDate, " End: ", endDate)
	return startDate, endDate
}

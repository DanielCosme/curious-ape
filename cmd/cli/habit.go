package main

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/api/types"
	"github.com/danielcosme/curious-ape/internal/client"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/go-sdk/colors"
	"github.com/danielcosme/go-sdk/dates"
	"github.com/spf13/cobra"
	"sort"
	"time"
)

var habitCmd = &cobra.Command{
	Use:               "habit <command> [flags]",
	Short:             "Manage habits",
	Long:              "List, create, update and delete habit logs",
	Aliases:           []string{"habits", "h"},
	ArgAliases:        []string{},
	PersistentPreRunE: loadCredentials,
	// Run:               listHabits,
}

var habitListCmd = &cobra.Command{
	Use:        "list",
	Short:      "habit list",
	Long:       `habit list`,
	Aliases:    listAliases,
	ArgAliases: []string{},
	Run:        listHabits,
}

func listHabits(cmd *cobra.Command, args []string) {
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
	keys := []time.Time{}
	for d, _ := range mapHabitsByDay {
		keys = append(keys, d)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Before(keys[j])
	})

	month := keys[0].Month().String()
	drawLine(fmt.Sprintf("%s (by %s)", colors.Blue(month), colors.Gray(by)), makeTableRow("S", "F", "W", "E"))
	drawLine("", "")
	for _, d := range keys {
		numberOfDay := d.Day()
		weekDay := d.Weekday()
		var sleep, fitness, work, eat entity.HabitStatus

		hs := mapHabitsByDay[d]
		for _, h := range hs {
			switch h.Type {
			case entity.HabitTypeWakeUp:
				sleep = h.Status
			case entity.HabitTypeFitness:
				fitness = h.Status
			case entity.HabitTypeDeepWork:
				work = h.Status
			case entity.HabitTypeFood:
				eat = h.Status
			}
		}

		var dot string
		if d.Before(dates.EndOfDay(time.Now())) {
			dot = IDot
		}

		habitsLine := makeTableRow(statusCheck(sleep), statusCheck(fitness), statusCheck(work), statusCheck(eat))
		drawLine(dayLine(numberOfDay, weekDay, dot), habitsLine)
	}
}

const IDot = "\uF444"
const SECTION_1_WIDTH = 40
const SECTION_2_WIDTH = 18

func statusCheck(status entity.HabitStatus) string {
	switch status {
	case entity.HabitStatusDone:
		return colors.Green("X")
	case entity.HabitStatusNotDone:
		return colors.Red("-")
	default:
		return " "
	}
}

func dayLine(numberOfDay int, weekOfTheDay time.Weekday, dot string) string {
	numDay := fmt.Sprintf("%-2d", numberOfDay)
	weekDay := fmt.Sprintf("%s", weekOfTheDay.String()[:2])
	d := colors.Purple(dot)
	dayLine := fmt.Sprintf("%s %s %s", colors.Yellow(numDay), weekDay, d)
	return dayLine
}

func drawLine(line1, line2 string) {
	fmt.Printf(" %-*s %-*s\n", SECTION_1_WIDTH, line1, SECTION_2_WIDTH, line2)
}

func makeTableRow(one, two, three, four string) string {
	return fmt.Sprintf("| %s | %s | %s | %s |", one, two, three, four)
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
		switch from {
		case "current":
			startDate = dates.StartOfTheMonth(time.Now())
			endDate = dates.EndOfTheMonth(time.Now())
		case "previous":
			startDate = dates.StartOfTheMonth(time.Now().AddDate(0, -1, 0))
			endDate = dates.EndOfTheMonth(time.Now().AddDate(0, -1, 0))
		}
	case "year":
	default:
		// All
	}
	return startDate, endDate
}

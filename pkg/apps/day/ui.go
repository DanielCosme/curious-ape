package day

import (
	"fmt"
	"net/url"
	"strconv"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	lucide "github.com/eduardolat/gomponents-lucide"

	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func UI_days(days []core.Day) Node {
	if len(days) == 0 {
		return Div(Text("No days available"))
	}

	next, prev := ui.GetNextPrevButtons(days[0], "")
	return Div(
		Class("days-container"),
		H2(Text(days[0].Date.Time().Month().String())),
		Div(
			Class("days-list"),
			Map(days, func(d core.Day) Node {
				return UI_day(d)
			}),
		),
		Div(
			Class("month-navigation"),
			next,
			prev,
		),
	)
}

func UI_day(day core.Day) Node {
	q := url.Values{}
	q.Add("date", day.Date.String())
	sync := fmt.Sprintf("@post('/day/sync?%s')", q.Encode())
	var goalIcon Node
	switch day.Habits.Score {
	case 0:
		goalIcon = lucide.Frown()
	case 1:
		goalIcon = lucide.HeartCrack()
	case 2:
		goalIcon = lucide.TriangleAlert()
	case 3:
		goalIcon = lucide.Trophy()
	case 4:
		goalIcon = lucide.ChessKing()
	default:
		goalIcon = nil
	}
	scoreClass := fmt.Sprintf("score-%d", day.Habits.Score)

	return Div(
		Class("day"),
		Span(Text(day.Date.Time().Format(core.HumanDate))),
		Span(
			Class("day-score flex "+scoreClass),
			goalIcon,
			Text(fmt.Sprintf("%d", day.Habits.Score)),
		),
		habitSpot(lucide.Bed(), day.Habits.Sleep),
		habitSpot(lucide.Dumbbell(), day.Habits.Fitness),
		habitSpot(lucide.UserCog(), day.Habits.DeepWork),
		habitSpot(lucide.Beef(), day.Habits.Eat),
		Button(Class(ui.CBtn+" btn-sync"), Text("sync"), ds.On("click", sync)),
		ID(fmt.Sprintf("day-%d", day.ID)),
	)
}

func habitSpot(icon Node, habit core.Habit) Node {
	q := url.Values{}
	q.Add("id", strconv.Itoa(int(habit.ID)))
	flipAction := fmt.Sprintf("@put('/habit/flip?%s')", q.Encode())

	var className string
	switch habit.State {
	case core.HabitStateDone:
		className = "habit-spot-done"
	case core.HabitStateNotDone:
		className = "habit-spot-not-done"
	default:
		className = ""
	}

	classes := "habit-spot"
	if className != "" {
		classes += " " + className
	}

	return Span(
		Class(classes),
		icon,
		If(habit.Note != "", Text(" "+habit.Note)),
		ds.On("click", flipAction),
	)
}

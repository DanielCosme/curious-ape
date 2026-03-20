package ui

import (
	"fmt"
	"net/url"
	"strconv"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
	lucide "github.com/eduardolat/gomponents-lucide"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

// Habit state symbols
const (
	habitDone    = "O"
	habitNotDone = "X"
	habitNoInfo  = "_"
)

func Home(s *State) Node {
	return bujoPage(s)
}

func bujoPage(s *State) Node {
	return layout("Days", s, days(s.Days))
}

func days(days []core.Day) Node {
	if len(days) == 0 {
		return Div(Text("No days available"))
	}

	next, prev := GetNextPrevButtons(days[0], "")
	return Div(
		Class("days-container"),
		H2(Text(days[0].Date.Time().Month().String())),
		Div(
			Class("days-list"),
			Map(days, func(d core.Day) Node {
				return Day(d)
			}),
		),
		Div(
			Class("month-navigation"),
			next,
			prev,
		),
	)
}

func Day(day core.Day) Node {
	q := url.Values{}
	q.Add("date", day.Date.String())
	sync := fmt.Sprintf("@post('/day/sync?%s')", q.Encode())
	var goalColor string
	var goalIcon Node
	switch day.Habits.Score {
	case 0:
		goalColor = "red"
		goalIcon = lucide.Frown()
	case 1:
		goalColor = "orangered"
		goalIcon = lucide.HeartCrack()
	case 2:
		goalColor = "goldenrod"
		goalIcon = lucide.TriangleAlert()
	case 3:
		goalColor = "yellowgreen"
		goalIcon = lucide.Trophy()
	case 4:
		goalColor = "rebeccapurple"
		goalIcon = lucide.ChessKing()
	default:
		goalColor = "white"
	}
	return Div(
		Class("day"),
		Span(Text(day.Date.Time().Format(core.HumanDate))),
		Span(
			Class("day-score flex"),
			Style(fmt.Sprintf("color: %s", goalColor)),
			goalIcon,
			Text(fmt.Sprintf("%d", day.Habits.Score)),
		),
		habitSpot(lucide.Bed(), day.Habits.Sleep),
		habitSpot(lucide.Dumbbell(), day.Habits.Fitness),
		habitSpot(lucide.UserCog(), day.Habits.DeepWork),
		habitSpot(lucide.Beef(), day.Habits.Eat),
		Button(Text("sync"), ds.On("click", sync)),
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

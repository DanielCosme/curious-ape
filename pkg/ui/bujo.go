package ui

import (
	"fmt"
	"net/url"
	"strconv"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
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
			Button(Text("Previous Month"), Disabled()),
			Button(Text("Next Month"), Disabled()),
		),
	)
}

func Day(day core.Day) Node {
	q := url.Values{}
	q.Add("date", day.Date.String())
	sync := fmt.Sprintf("@post('/day/sync?%s')", q.Encode())
	return Div(
		Class("day"),
		Span(Text(day.Date.Time().Format(core.HumanDate))),
		Span(Text("")),
		habitSpot(day.Habits.Sleep),
		habitSpot(day.Habits.Fitness),
		habitSpot(day.Habits.DeepWork),
		habitSpot(day.Habits.Eat),
		Button(Text("sync"), ds.On("click", sync)),
		ID(fmt.Sprintf("day-%d", day.ID)),
	)
}

func habitSpot(habit core.Habit) Node {
	state := habitNoInfo
	switch habit.State {
	case core.HabitStateDone:
		state = habitDone
	case core.HabitStateNotDone:
		state = habitNotDone
	}

	q := url.Values{}
	q.Add("id", strconv.Itoa(int(habit.ID)))
	flipAction := fmt.Sprintf("@put('/habit/flip?%s')", q.Encode())
	return Span(
		Class("habit-spot"),
		Strong(Text(state)),
		If(habit.Note != "", Text(" "+habit.Note)),
		ds.On("click", flipAction),
	)
}

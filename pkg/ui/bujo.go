package ui

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/danielcosme/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Home(s *State) Node {
	return bujoPage(s)
}

func bujoPage(s *State) Node {
	return layout("Days", s, days(s.Days))
}

func days(days []core.Day) Node {
	return Div(
		StyleEl(Text("span:hover { background-color: yellow; color: blue; cursor: pointer }")),
		H2(Text(days[0].Date.Time().Month().String())),
		Div(
			Map(days, func(d core.Day) Node {
				return Day(d)
			}),
		),
		Button(Text("Previous Month"), Disabled()),
		Button(Text("Next Month"), Disabled()),
	)
}

func Day(day core.Day) Node {
	q := url.Values{}
	q.Add("date", day.Date.String())
	sync := fmt.Sprintf("@post('/day/sync?%s')", q.Encode())
	return Div(
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
	state := "_"
	switch habit.State {
	case core.HabitStateDone:
		state = "O"
	case core.HabitStateNotDone:
		state = "X"
	}

	q := url.Values{}
	q.Add("id", strconv.Itoa(int(habit.ID)))
	flipAction := fmt.Sprintf("@put('/habit/flip?%s')", q.Encode())
	return Span(
		Strong(Text(state)),
		If(habit.Note != "", Text(" "+habit.Note)),
		ds.On("click", flipAction),
		Style("padding: 0 10px;"),
	)
}

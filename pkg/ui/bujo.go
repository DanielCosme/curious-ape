package ui

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

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

func GetNextPrevButtons(day core.Day, route string) (prev, next Node) {
	p, n := GetNextPrev(day, route)
	prev = Button(Text("Previous Month"), ds.On("click", p))
	next = Button(Text("Next Month"), ds.On("click", n))
	return
}

func GetNextPrev(day core.Day, route string) (prev, next string) {
	p, n := GetNextAndPreviousMonth(day)
	prev = fmt.Sprintf("@get('/%s?date=%s')", route, p)
	next = fmt.Sprintf("@get('/%s?date=%s')", route, n)
	return
}

func GetNextAndPreviousMonth(day core.Day) (prev, next string) {
	t := day.Date.FirstDayOfTheMonth().Time()
	previousMonth := t.AddDate(0, -1, 0)
	nextMonth := t.AddDate(0, 1, 0)
	now := time.Now()
	if previousMonth.Month() == now.Month() {
		previousMonth = now
	} else if nextMonth.Month() == now.Month() {
		nextMonth = now
	}
	prev = core.TimeFormatISO8601(previousMonth)
	next = core.TimeFormatISO8601(nextMonth)
	return
}

func Day(day core.Day) Node {
	q := url.Values{}
	q.Add("date", day.Date.String())
	sync := fmt.Sprintf("@post('/day/sync?%s')", q.Encode())
	return Div(
		Class("day"),
		Span(Text(day.Date.Time().Format(core.HumanDate))),
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

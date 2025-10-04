package views

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/danielcosme/curious-ape/pkg/core"

	. "github.com/delaneyj/gostar/elements"
)

type State struct {
	Version string
	Days    []core.Day
}

func Home(s *State) ElementRenderer {
	return bujoPage(s)
}

func Login(s *State) ElementRenderer {
	return layout(s, Group(
		H1().Text("Login"),
		FORM(
			DIV(
				LABEL(
					INPUT().TYPE("text").NAME("username").PLACEHOLDER(""),
				).Text("Username")),
			DIV(
				LABEL(
					INPUT().TYPE("text").NAME("password").PLACEHOLDER(""),
				).Text("Password")),
			BUTTON().Text("Login").
				TYPE("submit").
				DATASTAR_ON("click", "@post('/login', {contentType: 'form'})")),
	))
}

func bujoPage(s *State) ElementRenderer {
	return layout(s, days(s.Days))
}

func days(days []core.Day) ElementRenderer {
	return DIV(
		H2().Text(days[0].Date.Time().Month().String()),
		DIV(
			Range(days, func(d core.Day) ElementRenderer {
				return Day(d)
			}),
		),
	)
}

func Day(day core.Day) ElementRenderer {
	q := url.Values{}
	q.Add("date", day.Date.String())
	sync := fmt.Sprintf("@post('/day/sync?%s')", q.Encode())
	return DIV(
		SPAN().Text(day.Date.Time().Format(core.HumanDate)),
		SPAN().Text(""),
		habitSpot(day.Habits.Sleep),
		habitSpot(day.Habits.Fitness),
		habitSpot(day.Habits.DeepWork),
		habitSpot(day.Habits.Eat),
		BUTTON().Text("Sync").DATASTAR_ON("click", sync),
	).ID(fmt.Sprintf("day-%d", day.ID))
}

func habitSpot(habit core.Habit) ElementRenderer {
	state := "_"
	switch habit.State {
	case core.HabitStateDone:
		state = "O"
	case core.HabitStateNotDone:
		state = "X"
	}

	q := url.Values{}
	q.Add("id", strconv.Itoa(int(habit.ID)))
	flipAction := fmt.Sprintf("@put('/habit/flip?%s', {contentType: 'form'})", q.Encode())
	return SPAN(
		STYLE().Text("span:hover { background-color: yellow; color: blue; cursor: pointer }"),
		STRONG().Text(state),
	).DATASTAR_ON("click", flipAction).STYLE("padding", "0 10px")
}

func layout(s *State, children ...ElementRenderer) ElementRenderer {
	return Group(
		Text("<!DOCTYPE html>"),
		HEAD(
			META().CHARSET("UTF-8"),
			META().NAME("viewport").CONTENT("width=device-width, initial-scale=1.0"),
			SCRIPT().TYPE("module").SRC("https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.5/bundles/datastar.js"),
			TITLE().Text("Curious APE"),
		),
		HTML(
			BODY(
				HEADER(
					H1().Text("Curious Ape"),
				),
				NAV(
					a("/", "Home"),
					// a("/integrations", "Integrations"),
				),
				MAIN(children...),
				FOOTER(
					P().Text("Version "+s.Version),
				),
			),
		),
	)
}

func a(ref, name string) *AElement {
	return A().Text(name).HREF(ref)
}

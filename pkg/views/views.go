package views

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/pkg/core"

	. "github.com/delaneyj/gostar/elements"
)

type State struct {
	Version       string
	Authenticated bool
	Days          []core.Day
	Integrations  []application.IntegrationInfo
}

func Home(s *State) ElementRenderer {
	return bujoPage(s)
}

func DeepWork(s *State) ElementRenderer {
	p := Group(
		Range(s.Days, func(day core.Day) ElementRenderer {
			// TODO: Do this better.
			var duration time.Duration
			logs := []ElementRenderer{}
			for _, wl := range day.DeepWorkLogs {
				duration += wl.EndTime.Sub(wl.StartTime)
				logs = append(logs, SECTION(
					SPAN().Text(wl.Title),
					SPAN().Text(fmt.Sprintf("%s-%s", wl.StartTime.Format(core.Time), wl.EndTime.Format(core.Time))),
					SPAN().Text(fmt.Sprintf("  Duration: %s", core.DurationToString(wl.EndTime.Sub(wl.StartTime)))),
				))
			}
			return DIV(
				append([]ElementRenderer{H4().Text(day.Date.Time().Format(core.HumanDate)).Text("   " + core.DurationToString(duration))}, logs...)...,
			)
		}),
	)
	return layout(s, p)
}

func Fitness(s *State) ElementRenderer {
	p := Group(
		Range(s.Days, func(day core.Day) ElementRenderer {
			return DIV(
				Range(day.FitnessLogs, func(fl core.FitnessLog) ElementRenderer {
					return SECTION(
						H4().Text(fl.Title),
						SPAN().Text(fl.Date.Time().Format(core.HumanDate)+"                       "),
						SPAN().Text(fmt.Sprintf("%s-%s", fl.StartTime.Format(core.Time), fl.EndTime.Format(core.Time))),
						SPAN().Text(fmt.Sprintf("  Duration: %s", core.DurationToString(fl.EndTime.Sub(fl.StartTime)))),
					)
				}),
			)
		}),
	)
	return layout(s, p)
}

func Sleep(s *State) ElementRenderer {
	p := Group(
		Range(s.Days, func(day core.Day) ElementRenderer {
			return DIV(
				Range(day.SleepLogs, func(sl core.SleepLog) ElementRenderer {
					return SECTION(
						SPAN().Text(sl.Date.Time().Format(core.HumanDate)+"                       "),
						SPAN().Text(fmt.Sprintf("%s-%s", sl.StartTime.Format(core.Time), sl.EndTime.Format(core.Time))),
						SPAN().Text(fmt.Sprintf("  Duration: %s", core.DurationToString(sl.TimeAsleep))),
					)
				}),
			)
		}),
	)
	return layout(s, p)
}

func Login(s *State) ElementRenderer {
	return layout(s, Group(
		H1().Text("Login"),
		FORM(
			DIV(
				LABEL().Text("Username").Children(
					INPUT().TYPE("text").NAME("username").PLACEHOLDER("").STYLE("display", "block"),
				).STYLE("display", "block"),
			),
			DIV(
				LABEL().Text("Password").Children(
					INPUT().TYPE("text").NAME("password").PLACEHOLDER("").STYLE("display", "block"),
				).STYLE("display", "block"),
			),
			BUTTON().Text("Login").
				TYPE("submit").
				DATASTAR_ON("click", "@post('/login', {contentType: 'form'})")),
	))
}

func Integrations(s *State) ElementRenderer {
	return layout(s, DIV(
		Range(s.Integrations, func(i application.IntegrationInfo) ElementRenderer {
			return Integration(i)
		}),
	))
}

func Integration(i application.IntegrationInfo) ElementRenderer {
	integrationName := strings.ToLower(i.Name)
	q := url.Values{}
	q.Add("name", integrationName)
	onLoad := fmt.Sprintf("@get('/integration?%s')", q.Encode())

	return ARTICLE(
		H3().Text(i.Name),
		P().Text(fmt.Sprintf("Status: %s", i.Status)),
		If(len(i.Info) > 0,
			UL(
				Range(i.Info, func(info string) ElementRenderer {
					return LI().Text(info)
				}),
			),
		),
		If(i.Status == application.IntegrationStatusDicsonnected && i.AuthURL != "",
			A(
				BUTTON().Text("Authenticate"),
			).HREF(i.AuthURL).TARGET("_blank"),
		),
		// This component should do a GET requet to the server on load of the component of the dom.
	).ID("itg-"+integrationName).DATASTAR_ON("load", onLoad)
}

func bujoPage(s *State) ElementRenderer {
	return layout(s, days(s.Days))
}

func days(days []core.Day) ElementRenderer {
	return DIV(
		STYLE().Text("span:hover { background-color: yellow; color: blue; cursor: pointer }"),
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
	flipAction := fmt.Sprintf("@put('/habit/flip?%s')", q.Encode())
	return SPAN(
		STRONG().Text(state),
		If(habit.Note != "", Text(" "+habit.Note)),
	).DATASTAR_ON("click", flipAction).STYLE("padding", "0 10px")
}

func layout(s *State, children ...ElementRenderer) ElementRenderer {
	return Group(
		Text("<!DOCTYPE html>"),
		HTML(
			HEAD(
				META().CHARSET("UTF-8"),
				META().NAME("viewport").CONTENT("width=device-width, initial-scale=1.0"),
				SCRIPT().TYPE("module").SRC("https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.5/bundles/datastar.js"),
				LINK().REL("stylesheet").HREF("/assets/css/main.css"),
				TITLE().Text("Curious APE"),
			),
			BODY(
				HEADER(
					H1().Text("Curious Ape")),
				NAV().IfChildren(
					s.Authenticated,
					a("/", "Home"),
					a("/sleep", " Sleep "),
					a("/fitness", " Fitness "),
					a("/deep_work", " Deep-Work "),
					a("/integrations", "Integrations"),
					// TODO: Make pages loaded from nav be partially loaded, and not the full page.
				),
				MAIN(children...).IfChildren(
					s.Authenticated,
					BUTTON().Text("Logout").DATASTAR_ON("click", "@delete('/login')"),
				),
				FOOTER(
					P().Text(s.Version)),
			),
		),
	)
}

func a(ref, name string) *AElement {
	return A().Text(name).HREF(ref)
}

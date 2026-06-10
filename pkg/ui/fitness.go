package ui

import (
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"
	// ds "maragu.dev/gomponents-datastar"

	. "maragu.dev/gomponents/html"
)

func Fitness(s *State) Node {
	next, prev := GetNextPrevButtons(s.Days[0], "fitness")
	return layout("Fitness", s, Map(s.Days, func(day core.Day) Node {
		if len(day.FitnessLogs) == 0 {
			return nil
		}

		return Div(
			Map(day.FitnessLogs, func(fl core.FitnessLog) Node {
				return Section(
					H3(Text(fl.Title)),
					Span(Text(fl.Date.Time().Format(core.HumanDate))),
					Span(Text(fmt.Sprintf("%s-%s", fl.StartTime.Format(core.Time), fl.EndTime.Format(core.Time)))),
					Span(Text(fmt.Sprintf("  Duration: %s", core.DurationToString(fl.EndTime.Sub(fl.StartTime))))),
				)
			}),
		)
	}),
		Div(
			Class("month-navigation"),
			next,
			prev,
		),
	)
}

package fitnesslog

import (
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	. "maragu.dev/gomponents"

	// ds "maragu.dev/gomponents-datastar"

	. "maragu.dev/gomponents/html"
)

func UI_Fitness(s *ui.UIState, days []core.Day) Node {
	if len(days) == 0 {
		return ui.UILayout(s, P(Text("No records")))
	}

	next, prev := ui.GetNextPrevButtons(days[0].Date, "fitness")
	return ui.UILayout(s, Map(days, func(day core.Day) Node {
		if len(day.FitnessLogs) == 0 {
			return nil
		}

		return Div(
			Class(ui.CSurface),
			Map(day.FitnessLogs, func(fl core.FitnessLog) Node {
				return Div(
					Class(ui.CLogEntry),
					H3(Text(fl.Title)),
					Span(Text(fl.Date.Time.Format(core.HumanDate))),
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

package sleeplog

import (
	"fmt"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	. "maragu.dev/gomponents"

	// ds "maragu.dev/gomponents-datastar"

	. "maragu.dev/gomponents/html"
)

func UI_SleepLogPage(state *ui.State, days []core.Day) Node {
	state.Title = "Sleep"
	var content []Node
	if len(days) > 0 {
		next, prev := ui.GetNextPrevButtons(days[0].Date, "sleep")
		content = []Node{
			Map(days, func(day core.Day) Node {
				if len(day.SleepLogs) == 0 {
					return nil
				}
				return Div(
					Class(ui.CSurface),
					H3(Text(day.Date.Time.Format(core.HumanDate))),
					Map(day.SleepLogs, func(sl core.SleepLog) Node {
						return Div(
							Class(ui.CLogEntry),
							H4(Text(sl.Title)),
							P(Text(fmt.Sprintf("Wake up: %s", sl.EndTime.Format(core.Time)))),
							P(Text(fmt.Sprintf("  Duration: %s", core.DurationToString(sl.TimeAsleep)))),
						)
					}),
				)
			}),
			Div(
				Class("month-navigation"),
				next,
				prev,
			),
		}
	} else {
		content = []Node{
			P(Text("No records")),
		}
	}

	return ui.Layout(state, content...)

}

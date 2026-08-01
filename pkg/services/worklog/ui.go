package worklog

import (
	"fmt"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	. "maragu.dev/gomponents"

	// ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func UI_WorkLog(s *ui.UIState, days []core.Day) Node {
	next, prev := ui.GetNextPrevButtons(days[0].Date, "worklog")
	return ui.UILayout(s, Map(days, func(day core.Day) Node {
		if len(day.DeepWorkLogs) == 0 {
			return nil
		}
		var duration time.Duration
		nodes := []Node{}
		for _, wl := range day.DeepWorkLogs {
			duration += wl.EndTime.Sub(wl.StartTime)
			nodes = append(nodes, Div(
				Class(ui.CLogEntry),
				Span(Text(wl.Title+"  ")),
				Span(Text(fmt.Sprintf("%s-%s", wl.StartTime.Format(core.Time), wl.EndTime.Format(core.Time)))),
				Span(Text(fmt.Sprintf("  Duration: %s", core.DurationToString(wl.EndTime.Sub(wl.StartTime))))),
			))
		}
		return Div(
			Class(ui.CSurface),
			H3(Text(day.Date.Time().Format(core.HumanDate)+"   "+core.DurationToString(duration))),
			Group(nodes),
		)
	}),
		Div(
			Class("month-navigation"),
			next,
			prev,
		))
}

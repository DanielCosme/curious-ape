package ui

import (
	"fmt"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func DeepWork(s *State) Node {
	next, prev := GetNextPrevButtons(s.Days[0], "deep_work")
	return Layout("Deep Work", s, Map(s.Days, func(day core.Day) Node {
		if len(day.DeepWorkLogs) == 0 {
			return nil
		}
		var duration time.Duration
		nodes := []Node{}
		for _, wl := range day.DeepWorkLogs {
			duration += wl.EndTime.Sub(wl.StartTime)
			nodes = append(nodes, Div(
				Class(cLogEntry),
				Span(Text(wl.Title+"  ")),
				Span(Text(fmt.Sprintf("%s-%s", wl.StartTime.Format(core.Time), wl.EndTime.Format(core.Time)))),
				Span(Text(fmt.Sprintf("  Duration: %s", core.DurationToString(wl.EndTime.Sub(wl.StartTime))))),
			))
		}
		return Div(
			Class(cSurface),
			H3(Text(day.Date.Time().Format(core.HumanDate)+"   "+core.DurationToString(duration))),
			Group(nodes),
		)
	}),
		Div(
			Class("month-navigation"),
			next,
			prev,
		),
	)
}

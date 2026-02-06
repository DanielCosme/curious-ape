package ui

import (
	"fmt"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func DeepWork(s *State) Node {
	return layout("Deep Work", s, Map(s.Days, func(day core.Day) Node {
		var duration time.Duration
		nodes := []Node{}
		for _, wl := range day.DeepWorkLogs {
			duration += wl.EndTime.Sub(wl.StartTime)
			nodes = append(nodes, Section(
				Span(Text(wl.Title+"  ")),
				Span(Text(fmt.Sprintf("%s-%s", wl.StartTime.Format(core.Time), wl.EndTime.Format(core.Time)))),
				Span(Text(fmt.Sprintf("  Duration: %s", core.DurationToString(wl.EndTime.Sub(wl.StartTime))))),
			))
		}
		return Div(
			H4(Text(day.Date.Time().Format(core.HumanDate)+"   "+core.DurationToString(duration))),
			Group(nodes),
		)
	}))
}

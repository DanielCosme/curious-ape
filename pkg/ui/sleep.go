package ui

import (
	"fmt"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func Sleep(s *State) Node {
	return layout("Sleep", s, Map(s.Days, func(day core.Day) Node {
		if len(day.SleepLogs) == 0 {
			return nil
		}
		return Div(
			H3(Text(day.Date.Time().Format(core.HumanDate))),
			Map(day.SleepLogs, func(sl core.SleepLog) Node {
				return Section(
					H4(Text(sl.Title)),
					P(Text(fmt.Sprintf("Wake up: %s", sl.EndTime.Format(core.Time)))),
					P(Text(fmt.Sprintf("  Duration: %s", core.DurationToString(sl.TimeAsleep)))),
				)
			}),
		)
	}))
}

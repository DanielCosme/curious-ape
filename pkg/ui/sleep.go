package ui

import (
	"fmt"

	"github.com/danielcosme/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func Sleep(s *State) Node {
	return layout("Sleep", s, Map(s.Days, func(day core.Day) Node {
		return Div(
			Map(day.SleepLogs, func(sl core.SleepLog) Node {
				return Section(
					H4(Text(sl.Date.Time().Format(core.HumanDate))),
					P(Text(fmt.Sprintf("Wake up: %s", sl.EndTime.Format(core.Time)))),
					P(Text(fmt.Sprintf("  Duration: %s", core.DurationToString(sl.TimeAsleep)))),
				)
			}),
		)
	}))
}

package day

import (
	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/stephenafamo/bob"
)

var daysDB core.DayRepository
var habitsDB core.HabitRepository

func SetDaysBOB(db bob.DB) {
	daysDB = persistence.NewDays(db)
	habitsDB = persistence.NewHabits(db)
}

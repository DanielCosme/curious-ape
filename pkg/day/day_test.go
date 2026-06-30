package day_test

import (
	"testing"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/day"
	"danicos.dev/daniel/curious-ape/pkg/test"
	"github.com/stephenafamo/bob"
)

func TestDay(t *testing.T) {
	t.Parallel()
	bobDB := bob.NewDB(test.NewTestDB(t))
	app := day.New(bobDB)

	date1 := core.NewDate(time.Now()).FirstDayOfTheMonth()
	day1, err := app.GetOrCreate(date1)
	test.NilErr(t, err)
	test.True(t, day1.ID > 0)
	test.True(t, day1.Habits.Hs[0].State == core.HabitStateNoInfo)
	test.True(t, day1.Habits.Hs[0].Type == core.HabitTypeWakeUp)
	test.True(t, day1.Habits.Hs[1].State == core.HabitStateNoInfo)
	test.True(t, day1.Habits.Hs[1].Type == core.HabitTypeFitness)
	test.True(t, day1.Habits.Hs[2].State == core.HabitStateNoInfo)
	test.True(t, day1.Habits.Hs[2].Type == core.HabitTypeDeepWork)
	test.True(t, day1.Habits.Hs[3].State == core.HabitStateNoInfo)
	test.True(t, day1.Habits.Hs[3].Type == core.HabitTypeEatHealthy)

	date2 := core.NewDate(date1.Time().AddDate(0, 0, 1))
	days, err := app.Month(date2, core.ASC)
	test.NilErr(t, err)
	test.True(t, len(days) == 2)
	test.True(t, len(days[1].Habits.Hs) == 4)

	date3 := core.NewDate(date1.Time().AddDate(0, 0, 30))
	days, err = app.Month(date3, core.ASC)
	test.NilErr(t, err)
	test.True(t, len(days) == 31)
}

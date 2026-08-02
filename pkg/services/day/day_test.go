package day_test

import (
	"context"
	"testing"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	embbeded_nats "danicos.dev/daniel/curious-ape/pkg/mynats"
	"danicos.dev/daniel/curious-ape/pkg/services/day"
	"danicos.dev/daniel/curious-ape/pkg/services/habit"
	"danicos.dev/daniel/curious-ape/pkg/test"
	"github.com/stephenafamo/bob"
)

func TestDay(t *testing.T) {
	t.Parallel()

	t.Run("Month create expected range of dates", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ns, err := embbeded_nats.New(ctx)
		test.NilErr(t, err)
		ns.WaitForServer()
		nc, err := ns.Client()
		test.NilErr(t, err)

		bobDB := bob.NewDB(test.NewTestDB(t))
		srv := day.NewService(bobDB, nc)

		date1 := core.NewDate(time.Now()).FirstDayOfTheMonth()
		date2 := core.NewDate(date1.Time.AddDate(0, 0, 1))

		days, err := srv.Month(date2, core.ASC)
		test.NilErr(t, err)
		test.True(t, len(days) == 2)
		test.True(t, days[0].Date.IsEqual(date1.Time))
		test.True(t, days[1].Date.IsEqual(date2.Time))

		date3 := core.NewDate(date1.Time.AddDate(0, 0, 27))
		days, err = srv.Month(date3, core.ASC)
		test.NilErr(t, err)
		test.True(t, len(days) == 28)
	})

	t.Run("Day creation triggers habit creation, when habit service is initialized", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		ns, err := embbeded_nats.New(ctx)
		test.NilErr(t, err)
		ns.WaitForServer()
		nc, err := ns.Client()
		test.NilErr(t, err)

		bobDB := bob.NewDB(test.NewTestDB(t))
		_ = habit.NewService(bobDB, nc)
		srv := day.NewService(bobDB, nc)

		date1 := core.NewDate(time.Now()).FirstDayOfTheMonth()
		day1, err := srv.GetOrCreate(date1)
		test.NilErr(t, err)

		test.NilErr(t, err)
		test.True(t, day1.ID == 1)
		test.True(t, day1.Habits.Hs[0].State == core.HabitStateNoInfo)
		test.True(t, day1.Habits.Hs[0].Type == core.HabitTypeWakeUp)
		test.True(t, day1.Habits.Hs[1].State == core.HabitStateNoInfo)
		test.True(t, day1.Habits.Hs[1].Type == core.HabitTypeFitness)
		test.True(t, day1.Habits.Hs[2].State == core.HabitStateNoInfo)
		test.True(t, day1.Habits.Hs[2].Type == core.HabitTypeDeepWork)
		test.True(t, day1.Habits.Hs[3].State == core.HabitStateNoInfo)
		test.True(t, day1.Habits.Hs[3].Type == core.HabitTypeEatHealthy)
	})
}

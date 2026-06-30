package application_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"danicos.dev/daniel/curious-ape/pkg/application"
	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/oak"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"danicos.dev/daniel/curious-ape/pkg/test"
	"github.com/golang-migrate/migrate/v4"
	m_sqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/stephenafamo/bob"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func TestDeadline(t *testing.T) {
	t.Parallel()

	now := core.NewDateToday()
	testCases := map[string]struct {
		input    core.Deadline
		expected core.Deadline
		err      error
	}{
		"Succeeds with all fields": {
			input: core.Deadline{
				Title:     "Wife Birthday",
				StartDate: now,
				EndDate:   core.NewDate(now.Time().AddDate(0, 3, 0)),
				Recurring: true,
			},
			expected: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 1},
				Title:            "Wife Birthday",
				StartDate:        now,
				EndDate:          core.NewDate(now.Time().AddDate(0, 3, 0)),
				Recurring:        true,
			},
		},
		"Error when title is empty": {
			input: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "",
			},
			expected: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "",
			},
			err: errors.New("title is empty"),
		},
		"Error when start date is empty": {
			input: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "tt",
			},
			expected: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "tt",
			},
			err: errors.New("start time is empty"),
		},
		"Error when end date is empty": {
			input: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "tt",
				StartDate:        now,
			},
			expected: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "tt",
				StartDate:        now,
			},
			err: errors.New("end time is empty"),
		},
		"Error when end date is before start time": {
			input: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "tt",
				StartDate:        now,
				EndDate:          core.NewDate(now.Time().AddDate(0, -3, 0)),
			},
			expected: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "tt",
				StartDate:        now,
				EndDate:          core.NewDate(now.Time().AddDate(0, -3, 0)),
			},
			err: errors.New("end time cannot be before start time"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			app := NewTestApplication(t)
			ctx := context.Background()

			res, err := app.DeadlineCreate(ctx, tc.input)
			if tc.err == nil {
				test.NilErr(t, err)
			} else {
				test.True(t, strings.Contains(err.Error(), tc.err.Error()))
			}
			test.True(t, res.ID == tc.expected.ID)
			test.True(t, res.Title == tc.expected.Title)
			test.True(t, res.StartDate.Time().Equal(tc.expected.StartDate.Time()))
			test.True(t, res.EndDate.Time().Equal(tc.expected.EndDate.Time()))
			test.True(t, res.Recurring == tc.expected.Recurring)
		})
	}
}

/*

func TestDay(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)
	ctx := context.Background()

	date1 := core.NewDate(time.Now()).FirstDayOfTheMonth()
	day, err := app.DayGetOrCreate(date1)
	test.NilErr(t, err)
	test.True(t, day.ID > 0)

	params1 := core.Habit{Date: date1, Type: core.HabitTypeWakeUp, State: core.HabitStateDone}
	habit1, err := app.HabitUpsert(ctx, params1)
	test.NilErr(t, err)

	params2 := core.Habit{Date: date1, Type: core.HabitTypeFitness, State: core.HabitStateDone}
	habit2, err := app.HabitUpsert(ctx, params2)
	test.NilErr(t, err)

	params3 := core.Habit{Date: date1, Type: core.HabitTypeDeepWork, State: core.HabitStateDone}
	habit3, err := app.HabitUpsert(ctx, params3)
	test.NilErr(t, err)

	params4 := core.Habit{Date: date1, Type: core.HabitTypeEatHealthy, State: core.HabitStateDone}
	habit4, err := app.HabitUpsert(ctx, params4)
	test.NilErr(t, err)
	test.True(t, habit4.ID > 0)
	test.True(t, habit4.Date.IsEqual(date1.Time()))
	test.True(t, habit4.State == core.HabitStateDone)

	params4.State = core.HabitStateNotDone
	habit4, err = app.HabitUpsert(ctx, params4)
	test.NilErr(t, err)
	test.True(t, habit4.State == core.HabitStateNotDone)

	date2 := core.NewDate(date1.Time().AddDate(0, 0, 1))
	days, err := app.DaysMonth(ctx, date2)
	test.NilErr(t, err)
	test.True(t, len(days) == 2)
	test.True(t, len(days[0].Habits.Hs) == 4)
	test.True(t, days[1].Habits.Hs[0].Date.IsEqual(days[1].Date.Time()))
	test.True(t, days[1].Habits.Hs[0].ID == habit1.ID)
	test.True(t, days[1].Habits.Hs[1].ID == habit2.ID)
	test.True(t, days[1].Habits.Hs[2].ID == habit3.ID)
	test.True(t, days[1].Habits.Hs[3].ID == habit4.ID)
	test.True(t, len(days[0].Habits.Hs) == 4)
	test.True(t, days[0].Habits.Hs[0].State == core.HabitStateNoInfo)
	test.True(t, days[0].Habits.Hs[1].State == core.HabitStateNoInfo)
	test.True(t, days[0].Habits.Hs[2].State == core.HabitStateNoInfo)
	test.True(t, days[0].Habits.Hs[3].State == core.HabitStateNoInfo)

	date3 := core.NewDate(date1.Time().AddDate(0, 0, 30))
	days, err = app.DaysMonth(ctx, date3)
	test.NilErr(t, err)
	test.True(t, len(days) == 31)
}
*/

func TestApp_UserExists(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)

	exists, err := app.UserExists(1)
	test.NilErr(t, err)
	test.False(t, exists)

	err = app.SetPassword("daniel", "test", "admin@example.com", core.AuthRoleAdmin)
	test.NilErr(t, err)

	exists, err = app.UserExists(1)
	test.NilErr(t, err)
	test.True(t, exists)

	exists, err = app.UserExists(0)
	test.NilErr(t, err)
	test.False(t, exists)
}

func NewTestApplication(t *testing.T) *application.App {
	t.Helper()

	db := NewTestDB(t)
	t.Cleanup(func() { db.Close() })

	opts := &application.AppOptions{
		Logger: oak.New(oak.TintHandler(os.Stdout, oak.LevelTrace, true)),
		Config: &application.Config{
			Env: application.Test,
		},
		Database: persistence.New(bob.NewDB(db)),
	}
	app := application.New(opts)
	return app
}

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	failIfErr(t, err)

	migrationDriver, err := m_sqlite.WithInstance(db, &m_sqlite.Config{})
	failIfErr(t, err)

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://../../database/migrations/sqlite",
		"ape",
		migrationDriver,
	)
	failIfErr(t, err)

	err = migrator.Up()
	failIfErr(t, err)

	return db
}

func failIfErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

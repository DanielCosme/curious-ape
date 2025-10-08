package application_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/oak"
	"github.com/danielcosme/curious-ape/pkg/persistence"
	"github.com/danielcosme/curious-ape/pkg/test"
	"github.com/golang-migrate/migrate/v4"
	m_sqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/stephenafamo/bob"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

func TestHabitUpsertManual(t *testing.T) {
	t.Parallel()
	_ = NewTestApplication(t)
}

func TestDay(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)
	ctx := context.Background()

	date1 := core.NewDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
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
		Logger: oak.New(oak.TintHandler(os.Stdout, oak.LevelTrace)),
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

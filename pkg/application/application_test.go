package application_test

import (
	"database/sql"
	"log/slog"
	"testing"
	"time"

	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/fox"
	"github.com/danielcosme/curious-ape/pkg/persistence"
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

	date1 := core.NewDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	day, err := app.DayGetOrCreate(date1)
	fox.NilErr(t, err)
	fox.True(t, day.ID > 0)

	params1 := core.UpsertHabitParams{Date: date1, Type: core.HabitTypeWakeUp, State: core.HabitStateDone}
	habit1, err := app.HabitUpsert(params1)
	fox.NilErr(t, err)

	params2 := core.UpsertHabitParams{Date: date1, Type: core.HabitTypeFitness, State: core.HabitStateDone}
	habit2, err := app.HabitUpsert(params2)
	fox.NilErr(t, err)

	params3 := core.UpsertHabitParams{Date: date1, Type: core.HabitTypeDeepWork, State: core.HabitStateDone}
	habit3, err := app.HabitUpsert(params3)
	fox.NilErr(t, err)

	params4 := core.UpsertHabitParams{Date: date1, Type: core.HabitTypeEatHealthy, State: core.HabitStateDone}
	habit4, err := app.HabitUpsert(params4)
	fox.NilErr(t, err)
	fox.True(t, habit4.ID > 0)
	fox.True(t, habit4.Date.IsEqual(date1.Time()))
	fox.True(t, habit4.State == core.HabitStateDone)

	params4.State = core.HabitStateNotDone
	habit4, err = app.HabitUpsert(params4)
	fox.NilErr(t, err)
	fox.True(t, habit4.State == core.HabitStateNotDone)

	date2 := core.NewDate(date1.Time().AddDate(0, 0, 1))
	days, err := app.DaysMonth(date2)
	fox.NilErr(t, err)
	fox.True(t, len(days) == 2)
	fox.True(t, len(days[0].Habits) == 4)
	fox.True(t, days[0].Habits[0].Date.IsEqual(days[0].Date.Time()))
	fox.True(t, days[0].Habits[0].ID == habit1.ID)
	fox.True(t, days[0].Habits[1].ID == habit2.ID)
	fox.True(t, days[0].Habits[2].ID == habit3.ID)
	fox.True(t, days[0].Habits[3].ID == habit4.ID)
	fox.True(t, len(days[1].Habits) == 0)

	date3 := core.NewDate(date1.Time().AddDate(0, 0, 30))
	days, err = app.DaysMonth(date3)
	fox.NilErr(t, err)
	fox.True(t, len(days) == 31)
}

func TestApp_UserExists(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)

	exists, err := app.UserExists(1)
	fox.NilErr(t, err)
	fox.False(t, exists)

	err = app.SetPassword("daniel", "test", "admin@example.com", core.AuthRoleAdmin)
	fox.NilErr(t, err)

	exists, err = app.UserExists(1)
	fox.NilErr(t, err)
	fox.True(t, exists)

	exists, err = app.UserExists(0)
	fox.NilErr(t, err)
	fox.False(t, exists)
}

func NewTestApplication(t *testing.T) *application.App {
	t.Helper()

	db := NewTestDB(t)
	t.Cleanup(func() { db.Close() })

	opts := &application.AppOptions{
		Logger: slog.Default(),
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

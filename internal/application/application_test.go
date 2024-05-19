package application_test

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/internal/application"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database"
	"github.com/golang-migrate/migrate/v4"
	sqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/stephenafamo/bob"
	"gotest.tools/v3/assert"
	"log/slog"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func TestHabitUpsertManual(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)

	habit, err := app.HabitUpsert(core.NewHabitParams{
		Success:   true,
		Date:      core.NewDate(time.Now()),
		HabitType: core.HabitTypeWakeUp,
		Origin:    core.OriginLogSleep,
		Automated: true,
	})
	assert.NilError(t, err)
	assert.Assert(t, habit.IsZero() == false)
	assert.Assert(t, habit.State() == core.HabitStateDone)
	assert.Assert(t, len(habit.Logs) == 1)
	assert.DeepEqual(t, habit.Logs[0], core.HabitLog{
		ID:          1,
		Success:     true,
		IsAutomated: true,
		Origin:      core.OriginLogSleep,
	})

	habit, err = app.HabitUpsert(core.NewHabitParams{
		Success:   false,
		Date:      core.NewDate(time.Now()),
		HabitType: core.HabitTypeWakeUp,
		Origin:    core.OriginLogSleep,
		Automated: true,
	})
	assert.NilError(t, err)
	assert.Assert(t, habit.IsZero() == false)
	assert.Assert(t, habit.State() == core.HabitStateNotDone)
	assert.Assert(t, len(habit.Logs) == 1)

	habit, err = app.HabitUpsert(core.NewHabitParams{
		Success:   true,
		Date:      core.NewDate(time.Now()),
		HabitType: core.HabitTypeWakeUp,
		Origin:    core.OriginLogManual,
		Automated: false,
	})
	assert.NilError(t, err)
	assert.Assert(t, habit.IsZero() == false)
	assert.Assert(t, habit.State() == core.HabitStateDone)
	assert.Assert(t, len(habit.Logs) == 2)
	assert.Assert(t, habit.DayID == 1)
}

func TestDay(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)

	date1 := core.NewDate(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	day, err := app.DayGetOrCreate(date1)
	assert.NilError(t, err)
	assert.Assert(t, day.IsZero() == false)
	assert.Assert(t, day.ID > 0)

	_, err = app.HabitUpsert(core.NewHabitParams{Success: true, Date: day.Date, HabitType: core.HabitTypeWakeUp, Origin: core.OriginLogManual})
	assert.NilError(t, err)
	_, err = app.HabitUpsert(core.NewHabitParams{Success: true, Date: day.Date, HabitType: core.HabitTypeExercise, Origin: core.OriginLogManual})
	assert.NilError(t, err)
	_, err = app.HabitUpsert(core.NewHabitParams{Success: true, Date: day.Date, HabitType: core.HabitTypeDeepWork, Origin: core.OriginLogManual})
	assert.NilError(t, err)
	_, err = app.HabitUpsert(core.NewHabitParams{Success: true, Date: day.Date, HabitType: core.HabitTypeEatHealthy, Origin: core.OriginLogManual})
	assert.NilError(t, err)

	date2 := core.NewDate(date1.Time().AddDate(0, 0, 1))
	days, err := app.DaysMonth(date2)
	assert.NilError(t, err)
	assert.Assert(t, len(days) == 2)
	assert.Assert(t, len(days[0].Habits) == 4)
	assert.Assert(t, len(days[1].Habits) == 0)

	date3 := core.NewDate(date1.Time().AddDate(0, 0, 30))
	days, err = app.DaysMonth(date3)
	assert.NilError(t, err)
	assert.Equal(t, len(days), 31)
}

func TestApp_UserExists(t *testing.T) {
	t.Parallel()
	app := NewTestApplication(t)

	exists, err := app.UserExists(1)
	assert.NilError(t, err)
	assert.Assert(t, exists == false)

	err = app.SetPassword("daniel", "test", "admin@example.com", core.AuthRoleAdmin)
	assert.NilError(t, err)

	exists, err = app.UserExists(1)
	assert.NilError(t, err)
	assert.Assert(t, exists == true)

	exists, err = app.UserExists(0)
	assert.NilError(t, err)
	assert.Assert(t, exists == false)
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
		Database: database.New(bob.NewDB(db)),
	}
	app := application.New(opts)
	return app
}

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	failIfErr(t, err)

	migrationDriver, err := sqlite.WithInstance(db, &sqlite.Config{})
	failIfErr(t, err)

	migrator, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations/sqlite",
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

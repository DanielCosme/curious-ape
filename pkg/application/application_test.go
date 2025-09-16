package application_test

import (
	"database/sql"
	"github.com/danielcosme/curious-ape/pkg/application"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/persistence"
	"github.com/golang-migrate/migrate/v4"
	m_sqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/stephenafamo/bob"
	"gotest.tools/v3/assert"
	"log/slog"
	"testing"
	"time"

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
	assert.NilError(t, err)
	assert.Assert(t, day.ID > 0)

	_, err = app.HabitUpsert(date1, core.HabitTypeWakeUp, core.HabitStateDone)
	assert.NilError(t, err)
	_, err = app.HabitUpsert(date1, core.HabitTypeFitness, core.HabitStateDone)
	assert.NilError(t, err)
	_, err = app.HabitUpsert(date1, core.HabitTypeDeepWork, core.HabitStateDone)
	assert.NilError(t, err)
	_, err = app.HabitUpsert(date1, core.HabitTypeEatHealthy, core.HabitStateDone)
	assert.NilError(t, err)

	date2 := core.NewDate(date1.Time().AddDate(0, 0, 1))
	days, err := app.DaysMonth(date2)
	assert.NilError(t, err)
	assert.Assert(t, len(days) == 2)
	assert.Assert(t, len(days[0].R.Habits) == 4)
	assert.Assert(t, len(days[1].R.Habits) == 0)

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

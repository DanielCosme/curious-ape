package application_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"danicos.dev/daniel/curious-ape/pkg/application"
	"danicos.dev/daniel/curious-ape/pkg/config"
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
				EndDate:   core.NewDate(now.Time.AddDate(0, 3, 0)),
				Recurring: true,
			},
			expected: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 1},
				Title:            "Wife Birthday",
				StartDate:        now,
				EndDate:          core.NewDate(now.Time.AddDate(0, 3, 0)),
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
				EndDate:          core.NewDate(now.Time.AddDate(0, -3, 0)),
			},
			expected: core.Deadline{
				RepositoryCommon: core.RepositoryCommon{ID: 0},
				Title:            "tt",
				StartDate:        now,
				EndDate:          core.NewDate(now.Time.AddDate(0, -3, 0)),
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
			test.True(t, res.StartDate.Time.Equal(tc.expected.StartDate.Time))
			test.True(t, res.EndDate.Time.Equal(tc.expected.EndDate.Time))
			test.True(t, res.Recurring == tc.expected.Recurring)
		})
	}
}

func NewTestApplication(t *testing.T) *application.App {
	t.Helper()

	db := NewTestDB(t)
	t.Cleanup(func() { db.Close() })

	opts := &application.AppOptions{
		Logger: oak.New(oak.TintHandler(os.Stdout, oak.LevelTrace, true)),
		Config: &application.Config{
			Env: config.CI,
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

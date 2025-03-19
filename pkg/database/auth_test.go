package database

import (
	"database/sql"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/golang-migrate/migrate/v4"
	sqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stephenafamo/bob"
	"gotest.tools/v3/assert"
	"testing"
)

func TestAuthUpsert(t *testing.T) {
	db := NewTestDB(t)

	auth1, err := db.Auths.Upsert(&models.AuthSetter{
		Provider:     omit.From("fitbit"),
		AccessToken:  omit.From("access_token"),
		RefreshToken: omitnull.From("refresh_token"),
	})
	assert.NilError(t, err)
	assert.Assert(t, auth1 != nil)
	auth2, err := db.Auths.Upsert(&models.AuthSetter{
		Provider:     omit.From("fitbit"),
		AccessToken:  omit.From("access_token_2"),
		RefreshToken: omitnull.From("refresh_token_2"),
	})
	assert.NilError(t, err)
	assert.Assert(t, auth2 != nil)
	assert.Assert(t, auth1.ID == auth2.ID)
	assert.Assert(t, auth1.Provider == auth2.Provider)
}

func NewTestDB(t *testing.T) *Database {
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

	return New(bob.NewDB(db))
}

func failIfErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

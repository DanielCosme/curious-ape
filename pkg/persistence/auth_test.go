package persistence

import (
	"database/sql"
	"github.com/aarondl/opt/omit"
	"github.com/danielcosme/curious-ape/database/gen/models"
	"github.com/danielcosme/curious-ape/pkg/fox"
	"github.com/golang-migrate/migrate/v4"
	m_sqlite "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stephenafamo/bob"
	_ "modernc.org/sqlite"
	"testing"
)

func TestAuthUpsert(t *testing.T) {
	db := NewTestDB(t)

	auth1, err := db.Auths.Upsert(&models.OauthTokenSetter{
		Provider:     omit.From("fitbit"),
		AccessToken:  omit.From("access_token"),
		RefreshToken: omit.From("refresh_token"),
	})
	fox.NilErr(t, err)
	fox.False(t, auth1 == nil)
	auth2, err := db.Auths.Upsert(&models.OauthTokenSetter{
		Provider:     omit.From("fitbit"),
		AccessToken:  omit.From("access_token_2"),
		RefreshToken: omit.From("refresh_token_2"),
	})
	fox.NilErr(t, err)
	fox.False(t, auth2 == nil)
	fox.True(t, auth1.ID == auth2.ID)
	fox.True(t, auth1.Provider == auth2.Provider)
}

func NewTestDB(t *testing.T) *Database {
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

	return New(bob.NewDB(db))
}

func failIfErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

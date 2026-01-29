package main

import (
	"fmt"

	"github.com/danielcosme/curious-ape/pkg/config"
	"github.com/danielcosme/curious-ape/pkg/target"
	"github.com/magefile/mage/mg"
)

const (
	MigrateNew  = "new"
	MigrateUp   = "up"
	MigrateDown = "down"
)

type DB mg.Namespace

func (DB) Open() {
	t := target.NewA("sqlite3", "-box", dbLocation)
	r.RunV("db", t)
}

// Generates type-safe models to work with the database
func (DB) Gen() error {
	t := target.NewA("go", "tool", "bobgen-sqlite", "-c", "bobgen.yaml")
	return r.RunV("db:gen", t)
}

type Migrate mg.Namespace

// Creates new migration file takes 1 param <name>
func (Migrate) New(name string) error {
	migrate := target.New("migrate")
	migrate.Args("create", "-seq", "-ext=.sql", "-dir="+config.MIGRATIONS_LOCATION, name)
	return r.RunV("migrate:new", migrate)
}

func (Migrate) Up() error {
	return migrate("up", "")
}

func (Migrate) Down() error {
	return migrate("down", "")
}

// Forces specific schema_migration version
func (Migrate) Force(version string) error {
	return migrate("force", version)
}

func migrate(op, version string) error {
	t := target.New("migrate")
	dbConn := fmt.Sprintf("sqlite3://%s", dbLocation)
	t.Args("-verbose",
		"-path", config.MIGRATIONS_LOCATION,
		"-database", dbConn,
		op)
	if version != "" {
		t.Args(version)
	}
	return r.RunV("migrate:"+op, t)
}

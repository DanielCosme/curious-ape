package main

import (
	"database/sql"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
}

type RestAPI struct {
	App    *application.App
	Config *Config
}

func main() {
	db, err := sql.Open("sqlite3", "./ape.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	p := &RestAPI{
		App: application.New(db),
	}

	_, _ = p.App.Habits.Create(&entity.Habit{})
	fmt.Println("The end")
}

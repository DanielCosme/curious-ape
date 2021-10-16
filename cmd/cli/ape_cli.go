package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/danielcosme/curious-ape/internal/auth"
	"github.com/danielcosme/curious-ape/internal/data"
	"github.com/danielcosme/curious-ape/internal/validator"
)

const (
	layout   = "2006-01-02"
	credPath = "/.config/ape/cred.txt"
	baseUrl  = "https://ape.danicos.me/v1/"
	add      = "add"
	del      = "del"
	get      = "get"
	login    = "login"
	seed     = "seed"
	man      = "manual/cli"
)

// TODO installation script
type config struct {
	credentials string
}

type application struct {
	cfg config
}

func main() {
	var date string
	var app application

	today := time.Now().Format(layout)
	flag.StringVar(&date, "date", today, "date for habit to manipulate")
	flag.Parse()
	args := flag.Args()

	app.cfg.credentials = readCredentials()

	argsLen := len(args)
	if argsLen < 1 {
		log.Fatal("need at least 1 argument")
	}
	operation := args[0]

	switch operation {
	case add:
		if argsLen < 3 {
			fmt.Println("need at least 3 arguments")
			break
		}
		habit := &data.Habit{
			Type:   args[1],
			State:  args[2],
			Date:   date,
			Origin: man,
		}
		app.AddHabit(*habit)
	case del:
	case get:
	case login:
		usr := args[1]
		pas := args[2]
		app.Login(usr, pas)
	case seed:
		home := os.Getenv("HOME")
		filePath := home + "/.config/ape/habits.json"
		file, err := os.Create(filePath)
		panicIfErr(err)
		defer file.Close()

		var state string
		start, end := getDates()
		maxDay, maxMonth, maxYear := end.Date()
		habit := data.Habit{
			Type:   args[1],
			Origin: man,
		}

		currentDate := start
		for {
			curDay, curMonth, curYear := currentDate.Date()
			habit.Date = currentDate.Format(layout)
			fmt.Println("\nWhat is the", habit.Type, "habit state for:", habit.Date, "?")
			fmt.Scanln(&state)

			habit.State = state
			app.AddHabit(habit)

			currentDate = currentDate.AddDate(0, 0, 1)
			if maxDay == curDay && maxMonth == curMonth && maxYear == curYear {
				break
			}
		}

	default:
		fmt.Println("nope")
	}
}

func getDates() (first, last time.Time) {
	var start, end string
	fmt.Println("From which date?")
	fmt.Scanln(&start)
	first, err := time.Parse(layout, start)
	panicIfErr(err)

	fmt.Println("To which date?")
	fmt.Scanln(&end)
	last, err = time.Parse(layout, end)
	panicIfErr(err)

	if last.Before(first) {
		log.Fatal("last date cannot be smaller than first date")
	}
	return first, last
}

func (app *application) AddHabit(habit data.Habit) {
	fmt.Println("Adding", habit.Type, "habit for", habit.Date)
	validateHabit(&habit)

	jsonBody, err := json.Marshal(habit)
	panicIfErr(err)
	reader := strings.NewReader(string(jsonBody))

	url := baseUrl + "habits"
	res, err := app.makeRequest("POST", url, reader)
	panicIfErr(err)
	body, err := io.ReadAll(res.Body)
	panicIfErr(err)
	if res.StatusCode != http.StatusCreated {
		log.Fatal(res.Status, string(body))
	}
	fmt.Println("All Good!")
}

func validateHabit(habit *data.Habit) {
	val := validator.New()
	data.ValidateHabit(val, habit)
	if !val.Valid() {
		log.Fatal(val.Errors)
	}
}

func (app *application) Login(usr, pas string) {
	file, err := os.Create(getCredPath())
	panicIfErr(err)
	defer file.Close()

	fmt.Println("logging In...")
	encoded := auth.EncodeCredentials(usr, pas)
	app.cfg.credentials = encoded

	res, err := app.makeRequest("GET", baseUrl+"habits/1", nil)
	panicIfErr(err)
	if res.StatusCode == http.StatusUnauthorized {
		log.Fatal("invalid credentials")
	}

	fmt.Println("All Good!")
	file.WriteString(encoded)
}

func (app *application) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
	req, err := app.createRequest(method, url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (app *application) createRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", "Basic "+app.cfg.credentials)
	return req, err
}

func getCredPath() string {
	home := os.Getenv("HOME")
	return home + credPath
}

func openCredentials() *os.File {
	path := getCredPath()
	options := os.O_CREATE | os.O_RDONLY
	file, err := os.OpenFile(path, options, os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func readCredentials() string {
	file := openCredentials()
	defer file.Close()
	buf, err := io.ReadAll(file)
	panicIfErr(err)

	return string(buf)
}

func panicIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Clear() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

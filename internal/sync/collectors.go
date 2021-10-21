package sync

import (
	"github.com/danielcosme/curious-ape/internal/core"
	"log"
	"time"

	"github.com/danielcosme/curious-ape/internal/models"
	"github.com/danielcosme/curious-ape/internal/sync/fitbit"
	"github.com/danielcosme/curious-ape/internal/sync/google"
	"github.com/danielcosme/curious-ape/internal/sync/toggl"
)

type Collectors struct {
	Models *models.DB
	Sleep  *SleepCollector
	Work   *WorkCollector
	Fit    *FitnessCollector
}

func NewCollectors(models *models.DB) *Collectors {
	f := &SleepCollector{
		Models: models,
		SleepProvider: &fitbit.SleepProvider{
			Auth:  fitbit.FitbitAuth,
			Token: models.Tokens,
			Scope: "sleep",
		},
	}

	togg := &WorkCollector{
		Models: models,
		WorkProvider: &toggl.WorkProvider{
			Auth:  toggl.TogglAuth,
			Scope: "work",
		},
	}

	gGit := &FitnessCollector{
		Models: models,
		FitnessProvider: &google.FitnessProvider{
			Auth:  google.GoogleAuth,
			Token: models.Tokens,
			Scope: "fitness",
		},
	}

	return &Collectors{
		Models: models,
		Sleep:  f,
		Work:   togg,
		Fit:    gGit,
	}
}

func (co *Collectors) InitializeDayHabit() (err error) {
	t := time.Now().Format("2006-01-02")
	return co.InitializeDayHabits(t)
}

func (co *Collectors) InitializeDayHabits(date string) (err error) {
	types := []string{"sleep", "food", "fitness", "work"}
	h := core.Habit{
		State:  "no_info",
		Date:   date,
		Origin: "automated",
	}

	c := 0
	for _, v := range types {
		h.Type = v
		err = co.Models.Habits.Insert(&h)
		if err == nil {
			c++
		}
	}

	if err != nil {
		log.Println(c, "Habits Added,", err.Error())
		return err
	}

	log.Println(c, "CRON habits added successfully")
	return nil
}

func (co *Collectors) AllHabitsInit() error {
	start, _ := time.Parse("2006-01-02", "2021-01-01")
	end := time.Now().AddDate(0, 0, -1)
	maxDay, maxMonth, maxYear := end.Date()
	curDay, curMonth, curYear := start.Date()

	for {
		curDay, curMonth, curYear = start.Date()
		co.InitializeDayHabits(start.Format("2006-01-02"))
		start = start.AddDate(0, 0, 1)

		if maxDay == curDay && maxMonth == curMonth && maxYear == curYear {
			break
		}
	}

	return nil
}

package application_test

import (
	"fmt"
	application2 "github.com/danielcosme/curious-ape/internal/application"
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
	"os"
	"testing"
	"time"

	"github.com/danielcosme/curious-ape/internal/repository"
	logape "github.com/danielcosme/go-sdk/log"
	"gotest.tools/v3/assert"
)

func TestHabitUpsertManual(t *testing.T) {
	t.Parallel()

	app := NewTestApplication(t)
	date := entity2.NormalizeDate(time.Now(), time.UTC)
	data := &application2.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity2.HabitTypeFood),
		Success:      true,
		Origin:       entity2.Manual,
		IsAutomated:  false,
	}
	habit, err := app.HabitUpsert(data)
	assert.NilError(t, err)

	data = &application2.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity2.HabitTypeFood),
		Success:      false,
		Origin:       entity2.Manual,
		IsAutomated:  false,
	}

	habit, err = app.HabitUpsert(data)
	assert.NilError(t, err)

	for _, a := range habit.Logs {
		fmt.Println("s: ", a.Success)
	}
	{
		hs, err := app.HabitsGetAll(nil)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("In database")
		for _, h := range hs {
			for _, hl := range h.Logs {
				fmt.Println("s: ", hl.Success)
			}
		}
	}
	assert.DeepEqual(t, habit, &entity2.Habit{
		ID:         1,
		DayID:      1,
		CategoryID: 1,
		Status:     entity2.HabitStatusNotDone,
		Day: &entity2.Day{
			ID:   1,
			Date: date,
		},
		Category: &entity2.HabitCategory{
			ID:          1,
			Name:        "Eat healthy",
			Type:        "food",
			Code:        "food",
			Description: "",
			Color:       "#ffffff",
		},
		Logs: []*entity2.HabitLog{
			{
				ID:          1,
				HabitID:     1,
				Success:     false,
				Origin:      "manual",
				IsAutomated: false,
			},
		},
	})
}

func TestHabitUpsertAutomated(t *testing.T) {
	t.Parallel()

	app := NewTestApplication(t)

	date := entity2.NormalizeDate(time.Now(), time.UTC)
	data := &application2.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity2.HabitTypeFood),
		Success:      true,
		Origin:       entity2.Manual,
		Note:         "this is a test note",
		IsAutomated:  false,
	}

	habit, err := app.HabitUpsert(data)
	assert.NilError(t, err)

	data = &application2.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity2.HabitTypeFood),
		Success:      false,
		Origin:       entity2.Google,
		Note:         "automated entry",
		IsAutomated:  true,
	}

	habit, err = app.HabitUpsert(data)
	assert.NilError(t, err)

	assert.DeepEqual(t, habit, &entity2.Habit{
		ID:         1,
		DayID:      1,
		CategoryID: 1,
		Status:     entity2.HabitStatusDone,
		Day: &entity2.Day{
			ID:   1,
			Date: date,
		},
		Category: &entity2.HabitCategory{
			ID:          1,
			Name:        "Eat healthy",
			Type:        "food",
			Code:        "food",
			Description: "",
			Color:       "#ffffff",
		},
		Logs: []*entity2.HabitLog{
			{
				ID:          2,
				HabitID:     1,
				Success:     false,
				Note:        "automated entry",
				Origin:      "google",
				IsAutomated: true,
			},
			{
				ID:          1,
				HabitID:     1,
				Success:     true,
				Note:        "this is a test note",
				Origin:      "manual",
				IsAutomated: false,
			},
		},
	})

	data = &application2.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity2.HabitTypeWakeUp),
		Success:      false,
		Origin:       entity2.Fitbit,
		Note:         "automated entry from fitbit",
		IsAutomated:  true,
	}

	habit, err = app.HabitUpsert(data)
	assert.NilError(t, err)

	assert.DeepEqual(t, habit, &entity2.Habit{
		ID:         2,
		DayID:      1,
		CategoryID: 2,
		Status:     entity2.HabitStatusNotDone,
		Day: &entity2.Day{
			ID:   1,
			Date: date,
		},
		Category: &entity2.HabitCategory{
			ID:          2,
			Name:        "Wake up early",
			Type:        "wake_up",
			Code:        "wake_up",
			Description: "",
			Color:       "#ffffff",
		},
		Logs: []*entity2.HabitLog{
			{
				ID:          3,
				HabitID:     2,
				Success:     false,
				Note:        "automated entry from fitbit",
				Origin:      "fitbit",
				IsAutomated: true,
			},
		},
	})
}

func NewTestApplication(t *testing.T) *application2.App {
	t.Helper()

	// logger initialization
	logger := logape.New(os.Stdout, logape.LevelDebug, time.RFC3339)

	opts := &application2.AppOptions{
		Repository: repository.NewTestSqliteRepository(t),
		Logger:     logger,
		Config: &application2.Environment{
			Env: "test",
		},
	}

	app := application2.New(opts)
	return app
}

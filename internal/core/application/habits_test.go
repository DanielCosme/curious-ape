package application_test

import (
	"os"
	"testing"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/repository"
	logape "github.com/danielcosme/go-sdk/log"
	"gotest.tools/v3/assert"
)

func TestHabitUpsert(t *testing.T) {
	t.Parallel()

	app := NewTestApplication(t)

	date := entity.NormalizeDate(time.Now(), time.UTC)
	data := &application.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity.HabitTypeFood),
		Success:      true,
		Origin:       entity.Manual,
		Note:         "this is a test note",
		IsAutomated:  false,
	}

	habit, err := app.HabitUpsert(data)
	assert.NilError(t, err)

	data = &application.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity.HabitTypeFood),
		Success:      false,
		Origin:       entity.Google,
		Note:         "automated entry",
		IsAutomated:  true,
	}

	habit, err = app.HabitUpsert(data)
	assert.NilError(t, err)

	assert.DeepEqual(t, habit, &entity.Habit{
		ID:         1,
		DayID:      1,
		CategoryID: 1,
		Status:     entity.HabitStatusDone,
		Day: &entity.Day{
			ID:   1,
			Date: date,
		},
		Category: &entity.HabitCategory{
			ID:          1,
			Name:        "Eat healthy",
			Type:        "food",
			Code:        "food",
			Description: "",
			Color:       "#ffffff",
		},
		Logs: []*entity.HabitLog{
			{
				ID:          1,
				HabitID:     1,
				Success:     true,
				Note:        "this is a test note",
				Origin:      "manual",
				IsAutomated: false,
			},
			{
				ID:          2,
				HabitID:     1,
				Success:     false,
				Note:        "automated entry",
				Origin:      "google",
				IsAutomated: true,
			},
		},
	})

	data = &application.NewHabitParams{
		Date:         time.Now(),
		CategoryCode: string(entity.HabitTypeWakeUp),
		Success:      false,
		Origin:       entity.Fitbit,
		Note:         "automated entry from fitbit",
		IsAutomated:  true,
	}

	habit, err = app.HabitUpsert(data)
	assert.NilError(t, err)

	assert.DeepEqual(t, habit, &entity.Habit{
		ID:         2,
		DayID:      1,
		CategoryID: 2,
		Status:     entity.HabitStatusNotDone,
		Day: &entity.Day{
			ID:   1,
			Date: date,
		},
		Category: &entity.HabitCategory{
			ID:          2,
			Name:        "Wake up early",
			Type:        "wake_up",
			Code:        "wake_up",
			Description: "",
			Color:       "#ffffff",
		},
		Logs: []*entity.HabitLog{
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

func NewTestApplication(t *testing.T) *application.App {
	t.Helper()

	// logger initialization
	logger := logape.New(os.Stdout, logape.LevelTrace, time.RFC3339)
	logape.DefaultLogger = logger

	opts := &application.AppOptions{
		Repository: repository.NewTestSqliteRepository(t),
		Logger:     logger,
		Config: &application.Environment{
			Env: "test",
		},
		SessionManager: nil,
	}

	app := application.New(opts)
	return app
}

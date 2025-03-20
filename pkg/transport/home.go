package transport

import (
	"fmt"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DaysPayload struct {
	Month string       `json:"month"`
	Days  []DaySummary `json:"days"`
}

type DaySummary struct {
	Key     string       `json:"key"`
	Date    string       `json:"date"`
	Day     string       `json:"day"`
	WakeUp  HabitSummary `json:"wake_up"`
	Fitness HabitSummary `json:"fitness"`
	Work    HabitSummary `json:"work"`
	Eat     HabitSummary `json:"eat"`
}

type HabitSummary struct {
	State core.HabitState `json:"state"`
	Type  core.HabitType  `json:"type"`
}

func (t *Transport) home(c echo.Context) error {
	days, err := t.App.DaysMonth(core.NewDateToday())
	if err != nil {
		return errServer(err)
	}
	daysPayload := DaysPayload{
		Month: days[0].Date.Month().String(),
	}
	for _, day := range days {
		daysPayload.Days = append(daysPayload.Days, dayDBToSummary(day))
	}
	return c.JSON(http.StatusOK, daysPayload)
}

func dayDBToSummary(day *models.Day) DaySummary {
	format := core.TimeFormatISO8601(day.Date)
	ds := DaySummary{
		Key:     fmt.Sprintf("day_s_%s", format),
		Date:    format,
		Day:     day.Date.Format(core.HumanDate),
		WakeUp:  HabitSummary{State: core.HabitStateNoInfo, Type: core.HabitTypeWakeUp},
		Fitness: HabitSummary{State: core.HabitStateNoInfo, Type: core.HabitTypeFitness},
		Work:    HabitSummary{State: core.HabitStateNoInfo, Type: core.HabitTypeDeepWork},
		Eat:     HabitSummary{State: core.HabitStateNoInfo, Type: core.HabitTypeEatHealthy},
	}
	for _, h := range day.R.Habits {
		switch core.HabitType(h.R.HabitCategory.Kind) {
		case core.HabitTypeWakeUp:
			ds.WakeUp = habitDBToTransport(h)
		case core.HabitTypeFitness:
			ds.Fitness = habitDBToTransport(h)
		case core.HabitTypeDeepWork:
			ds.Work = habitDBToTransport(h)
		case core.HabitTypeEatHealthy:
			ds.Eat = habitDBToTransport(h)
		}
	}
	return ds
}

func habitDBToTransport(h *models.Habit) HabitSummary {
	return HabitSummary{
		Type:  core.HabitType(h.R.HabitCategory.Kind),
		State: core.HabitState(h.State),
	}
}

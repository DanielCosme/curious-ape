package transport

import (
	"fmt"
	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/database/gen/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"sort"
	"time"
)

type DaysPayload struct {
	Month string       `json:"month"`
	Days  []DaySummary `json:"days"`
}

type DaySummary struct {
	Key           string       `json:"key"`
	Date          string       `json:"date"`
	Day           string       `json:"day"`
	WakeUp        HabitSummary `json:"wake_up_habit"`
	Fitness       HabitSummary `json:"fitness_habit"`
	Work          HabitSummary `json:"work_habit"`
	Eat           HabitSummary `json:"eat_habit"`
	WakeUpDetail  string       `json:"wake_up_detail"`
	FitnessDetail string       `json:"fitness_detail"`
	WorkDetail    string       `json:"work_detail"`
	Score         int          `json:"score"`
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
	daysPayload := DaysPayload{Month: days[0].Date.Month().String()}
	sort.Sort(DaysSliceSort(days))
	for _, day := range days {
		daysPayload.Days = append(daysPayload.Days, dayDBToSummary(day))
	}
	return c.JSON(http.StatusOK, daysPayload)
}

func (t *Transport) syncDay(c echo.Context) error {
	day, err := core.DateFromISO8601(c.QueryParam("day"))
	if err != nil {
		return errClientError(fmt.Errorf("invalid date param - %w", err))
	}
	dayDB, err := t.App.SyncDay(day)
	if err != nil {
		return errServer(err)
	}
	return c.JSON(http.StatusOK, dayDBToSummary(dayDB))
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
		if h.State == string(core.HabitStateDone) {
			ds.Score++
		}
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
	for _, sl := range day.R.SleepLogs {
		if sl.IsMainSleep.GetOrZero() {
			ds.WakeUpDetail = sl.EndTime.Format(core.Time)
			break
		}
	}
	for idx, f := range day.R.FitnessLogs {
		// 15:21 - 16:05 (44m0s)
		ds.FitnessDetail = fmt.Sprintf(
			"%s - %s (%s)",
			f.StartTime.Format(core.Time),
			f.EndTime.Format(core.Time),
			f.EndTime.Sub(f.StartTime).Round(time.Minute))
		// TODO: figure out how to handle this once I have more fitness logs per day.
		if idx == 0 {
			break
		}
	}
	var workDuration time.Duration
	for _, wl := range day.R.DeepWorkLogs {
		workDuration += time.Duration(wl.Seconds) * time.Second
	}
	ds.WorkDetail = workDuration.Round(time.Minute).String()
	return ds
}

func habitDBToTransport(h *models.Habit) HabitSummary {
	return HabitSummary{
		Type:  core.HabitType(h.R.HabitCategory.Kind),
		State: core.HabitState(h.State),
	}
}

type DaysSliceSort []*models.Day

func (a DaysSliceSort) Len() int           { return len(a) }
func (a DaysSliceSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DaysSliceSort) Less(i, j int) bool { return a[i].Date.After(a[j].Date) }

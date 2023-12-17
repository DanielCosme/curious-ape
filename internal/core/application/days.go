package application

import (
	"fmt"
	"time"

	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/curious-ape/internal/integrations/toggl"
	"github.com/danielcosme/go-sdk/errors"
	"github.com/danielcosme/go-sdk/log"
)

func (a *App) DaysGetAll() ([]*entity.Day, error) {
	return a.db.Days.Find(entity.DayFilter{}, database.DaysPipeline(a.db)...)
}

func (a *App) DayGetByID(id int) (*entity.Day, error) {
	return a.db.Days.Get(entity.DayFilter{IDs: []int{id}}, database.DaysJoinHabits(a.db))
}

func (a *App) SyncDeepWorkByDateRange(start, end time.Time) error {
	togglAPI, o, err := a.TogglAPI()
	if err != nil {
		return err
	}
	days, err := a.daysGetByDateRange(start, end)
	if err != nil {
		return err
	}

	for _, d := range days {
		summary, err := togglAPI.Reports.GetDaySummaryForProjectIDs(d.Date, o.ToogglProjectIDs, o.ToogglWorkSpaceID)
		if err != nil {
			return err
		}
		if _, err := a.dayUpdate(d, workLogFromToggl(summary)); err != nil {
			return nil
		}

		if err := a.createDeepWorkLog(d, entity.Toggl); err != nil {
			return err
		}
		togglSleep()
	}

	return nil
}

func (a *App) SyncDeepWorkLog(date time.Time) error {
	togglAPI, o, err := a.TogglAPI()
	if err != nil {
		return err
	}
	d, err := database.DayGetOrCreate(a.db, date)
	if err != nil {
		return err
	}

	summary, err := togglAPI.Reports.GetDaySummaryForProjectIDs(d.Date, o.ToogglProjectIDs, o.ToogglWorkSpaceID)
	if err != nil {
		return err
	}
	if _, err := a.dayUpdate(d, workLogFromToggl(summary)); err != nil {
		return nil
	}

	return a.createDeepWorkLog(d, entity.Toggl)
}

func (a *App) TogglAPI() (*toggl.API, *entity.Auth, error) {
	o, err := a.db.Auths.Get(entity.AuthFilter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
	if err != nil {
		return nil, nil, err
	}
	return a.sync.TogglClient(o.AccessToken), o, nil
}

func (a *App) SyncDeepWork() error {
	days, err := a.db.Days.Find(entity.DayFilter{}, database.DaysJoinSleepLogs(a.db))
	if err != nil {
		return err
	}
	togglAPI, o, err := a.TogglAPI()
	if err != nil {
		return err
	}

	for _, d := range days {
		summary, err := togglAPI.Reports.GetDaySummaryForProjectIDs(d.Date, o.ToogglProjectIDs, o.ToogglWorkSpaceID)
		if err != nil {
			return err
		}
		if _, err := a.dayUpdate(d, workLogFromToggl(summary)); err != nil {
			return nil
		}

		if err := a.createDeepWorkLog(d, entity.Toggl); err != nil {
			return err
		}
		togglSleep()
	}

	return nil
}

func workLogFromToggl(s *toggl.Summary) *entity.Day {
	return &entity.Day{
		DeepWorkMinutes: int(toggl.ToDuration(s.TotalGrand).Minutes()),
	}
}

func togglSleep() {
	// Toggle Api Only accepts 1 api cal per second
	time.Sleep(time.Second)
}

func (a *App) HabitUpsertFromDeepWorkLog(d *entity.Day, origin entity.DataSource) error {
	habitCategory, err := a.HabitCategoryGetByType(entity.HabitTypeDeepWork)
	if err != nil {
		return err
	}

	var success bool
	// If the deep work minutes are bigger than 5 hours
	dur := time.Duration(d.DeepWorkMinutes) * time.Minute
	if dur >= (time.Hour * 5) {
		success = true
	}

	_, err = a.HabitUpsert(&NewHabitParams{
		Date:         d.Date,
		CategoryCode: habitCategory.Code,
		Success:      success,
		Origin:       origin,
		Note:         fmt.Sprintf("Deep work duration: %s", dur.String()),
		IsAutomated:  false,
	})
	return err
}

func (a *App) DayUpdate(day, data *entity.Day) (*entity.Day, error) {
	var err error
	day, err = a.dayUpdate(day, data)
	if err != nil {
		return nil, err
	}
	// create deep work resource (in the future) and upsert habit
	if err := a.createDeepWorkLog(day, entity.Manual); err != nil {
		return nil, err
	}
	return day, database.ExecuteDaysPipeline([]*entity.Day{day}, database.DaysJoinHabits(a.db))
}

func (a *App) dayUpdate(day, data *entity.Day) (*entity.Day, error) {
	day.DeepWorkMinutes = data.DeepWorkMinutes
	return a.db.Days.Update(day, database.DaysPipeline(a.db)...)
}

func (a *App) createDeepWorkLog(day *entity.Day, origin entity.DataSource) error {
	if err := a.HabitUpsertFromDeepWorkLog(day, origin); err != nil {
		return err
	}

	a.Log.InfoP("updated deep work log", log.Prop{
		"origin": origin.Str(),
		"date":   day.Date.Format(entity.HumanDateWithTime),
	})
	return nil
}

func (a *App) daysGetByDateRange(start, end time.Time) ([]*entity.Day, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("invalid dates")
	}
	if start.After(end) {
		return nil, errors.New("start date must be before end")
	}
	return a.db.Days.Find(entity.DayFilter{Dates: datesRange(start, end)})
}

func (a *App) DaysMonth() ([]*entity.Day, error) {
	var res []*entity.Day
	var err error
	var d *entity.Day

	today := time.Now()
	first := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	d, err = a.db.Days.Get(entity.DayFilter{Dates: []time.Time{today}})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}
	if d != nil {
		var e error
		res, e = a.daysGetByDateRange(first, today)
		if e != nil {
			return nil, err
		}
		return res, database.ExecuteDaysPipeline(res, database.DaysPipeline(a.db)...)
	}

	for _, dt := range datesRange(first, today) {
		d, err = database.DayGetOrCreate(a.db, dt)
		if err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, database.ExecuteDaysPipeline(res, database.DaysJoinHabits(a.db))
}

func datesRange(start, end time.Time) []time.Time {
	dates := []time.Time{}
	inter := start

	for inter.Before(end) {
		dates = append(dates, inter)
		inter = inter.AddDate(0, 0, 1)
	}
	dates = append(dates, end)

	return dates
}

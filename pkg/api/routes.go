package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/dove"
	"github.com/danielcosme/curious-ape/pkg/views"
)

func Routes(a *API) http.Handler {
	d := dove.New(a.App.Log.Handler())

	// e.StaticFS("/static", echo.MustSubFS(views.StaticFS, "static"))

	d.Use(dove.MiddlewarePanicRecover)

	d.Endpoint("/").GET(a.Home)
	d.Endpoint("/habit/flip").PUT(a.HabitFlip)
	d.Endpoint("/day/sync").POST(a.DaySync)

	return d
}

func (a *API) Home(c *dove.Context) error {
	days, err := a.App.DaysMonth(c.Ctx(), core.NewDate(time.Now()))
	if err != nil {
		return err
	}

	s := &views.State{Days: days}
	return c.RenderOK(views.Home(s))
}

func (a *API) HabitFlip(c *dove.Context) error {
	c.ParseForm()
	id, _ := strconv.Atoi(c.Req.Form.Get("id"))
	habit, err := a.App.HabitFlip(id)
	if err != nil {
		return err
	}
	day, err := a.App.DayGetOrCreate(habit.Date)
	if err != nil {
		return err
	}
	return c.RenderOK(views.Day(day))
}

func (a *API) DaySync(c *dove.Context) error {
	c.ParseForm()
	date, _ := core.DateFromISO8601(c.Req.Form.Get("date"))
	day, err := a.App.DaySync(c.Ctx(), date)
	if err != nil {
		return err
	}
	return c.RenderOK(views.Day(day))
}

func (a *API) State() views.State {
	return views.State{}
}

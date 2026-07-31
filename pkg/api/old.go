package api

import (
	"net/http"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/config"
	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/dove"
	"danicos.dev/daniel/curious-ape/pkg/ui"
)

func Routes(a *OldAPI) http.Handler {
	d := dove.New(a.App.Log.Handler())

	d.Use(dove.MiddlewarePanicRecover)

	if a.App.Env == config.Dev {
		d.Use(DevMiddleware)
	}

	d.Use(a.MiddlewareLoadCookie)
	d.Use(a.MiddlewareAuthenticateFromSession)
	d.Use(a.MiddlewareSetUIState)

	d.Endpoint("/api/oauth2/fitbit/success").GET(a.FitbitSuccess)
	d.Endpoint("/api/oauth2/google/success").GET(a.GoogleSuccess)

	d.Use(a.MiddlewareRequireAuthentication)

	d.Endpoint("/day/sync").POST(a.DaySync)
	d.Endpoint("/sleep").GET(a.Sleep)
	d.Endpoint("/fitness").GET(a.Fitness)
	d.Endpoint("/deep_work").GET(a.DeepWork)
	d.Endpoint("/deadlines").GET(a.DeadlinesList)
	d.Endpoint("/deadline").
		GET(a.DeadlinesGetForm).
		POST(a.DeadlinesPostForm)

	return d
}

func (a *OldAPI) DeadlinesList(c *dove.Context) error {
	res, err := a.App.DeadlineList(c.Ctx())
	if err == nil {
		state := State(a, c.Req)
		state.Deadlines.DS = res
		return c.RenderOK(ui.Deadlines(state))
	}
	return err
}

func (a *OldAPI) DeadlinesGetForm(c *dove.Context) error {
	state := State(a, c.Req)
	return c.RenderOK(ui.DeadlineForm(state))
}

func (a *OldAPI) DeadlinesPostForm(c *dove.Context) error {
	c.ParseForm()
	state := State(a, c.Req)
	var recurring bool
	if c.Req.PostForm.Get("recurrent") == "on" {
		recurring = true
	}
	date, err := core.NewDateFromISO8601(c.Req.PostForm.Get("end_date"))
	if err == nil {
		_, err := a.App.DeadlineCreate(c.Ctx(), core.Deadline{
			Title:     c.Req.PostForm.Get("title"),
			StartDate: core.NewDateToday(),
			EndDate:   date,
			Recurring: recurring,
		})
		if err == nil {
			return c.Redirect("/deadlines")
		}
		state.Deadlines.Err = err
		return c.RenderOK(ui.DeadlineForm(state))
	}
	return err
}

func (a *OldAPI) DeepWork(c *dove.Context) error {
	days, err := a.App.Day.Month(getDateParam(c), core.DESC)
	if err == nil {
		state := State(a, c.Req)
		state.Days = days
		return c.RenderOK(ui.DeepWork(state))
	}
	return err
}

func (a *OldAPI) Fitness(c *dove.Context) error {
	days, err := a.App.Day.Month(getDateParam(c), core.DESC)
	if err == nil {
		state := State(a, c.Req)
		state.Days = days
		return c.RenderOK(ui.Fitness(state))
	}
	return err
}

func (a *OldAPI) Sleep(c *dove.Context) error {
	days, err := a.App.Day.Month(getDateParam(c), core.DESC)
	if err == nil {
		state := State(a, c.Req)
		state.Days = days
		return c.RenderOK(ui.Sleep(state))
	}
	return err
}

func (a *OldAPI) FitbitSuccess(c *dove.Context) error {
	c.ParseForm()
	return a.App.Oauth2Success(core.IntegrationFitbit, c.Req.FormValue("code"))
}

func (a *OldAPI) GoogleSuccess(c *dove.Context) error {
	c.ParseForm()
	return a.App.Oauth2Success(core.IntegrationGoogle, c.Req.FormValue("code"))
}

func getDateParam(c *dove.Context) core.Date {
	c.ParseForm()
	if c.Req.Form.Get("date") == "" {
		return core.NewDate(time.Now())
	} else {
		date, err := core.NewDateFromISO8601(c.Req.Form.Get("date"))
		if err == nil {
			return date
		}
		c.Log.Fatal("cannot parse date", "err", err)
		panic(err)
	}
}

func (a *OldAPI) DaySync(c *dove.Context) error {
	c.ParseForm()
	date, _ := core.NewDateFromISO8601(c.Req.Form.Get("date"))
	_, err := a.App.DaySync(c.Ctx(), date)
	if err == nil {
		// return c.RenderOK(day.UI_day(d))
	}
	return err
}

func State(a *OldAPI, r *http.Request) *ui.State {
	return &ui.State{
		Version:       a.Version,
		Authenticated: a.IsAuthenticated(r),
		CurrentPath:   r.URL.Path,
	}
}

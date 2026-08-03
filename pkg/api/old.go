package api

import (
	"net/http"

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

	d.Use(a.MiddlewareRequireAuthentication)

	d.Endpoint("/deadline").
		GET(a.DeadlinesGetForm).
		POST(a.DeadlinesPostForm)

	return d
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

func State(a *OldAPI, r *http.Request) *ui.State {
	return &ui.State{
		Version:       a.Version,
		Authenticated: a.IsAuthenticated(r),
		CurrentPath:   r.URL.Path,
	}
}

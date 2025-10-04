package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/dove"
	"github.com/danielcosme/curious-ape/pkg/oak"
	"github.com/danielcosme/curious-ape/pkg/persistence"
	"github.com/danielcosme/curious-ape/pkg/views"
)

func Routes(a *API) http.Handler {
	d := dove.New(a.App.Log.Handler())

	// e.StaticFS("/static", echo.MustSubFS(views.StaticFS, "static"))

	d.Use(dove.MiddlewarePanicRecover)
	d.Use(a.MiddlewareLoadCookie)
	d.Use(a.MiddlewareAuthenticateFromSession)

	d.Endpoint("/login").
		GET(a.GetLoginForm).
		POST(a.Login).
		DELETE(a.Logout)

	d.Use(a.MiddlewareRequireAuthentication)

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

	s := State(a, c.Req)
	s.Days = days
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

func (a *API) GetLoginForm(c *dove.Context) error {
	if a.IsAuthenticated(c.Req) {
		return c.Redirect("/")
	}
	return c.RenderOK(views.Login(State(a, c.Req)))
}

func (a *API) Login(c *dove.Context) error {
	logger := oak.FromContext(c.Ctx())

	c.ParseForm()
	username := c.Req.PostFormValue("username")
	password := c.Req.PostFormValue("password")
	id, err := a.App.Authenticate(username, password)
	if err != nil {
		if errors.Is(err, persistence.ErrInvalidCredentials) {
			// TODO: send http.StatusUnauthorized
			return err
		} else {
			// TODO: send http.InternalServerError
			return err
		}
	}
	err = a.Scs.RenewToken(c.Ctx())
	if err != nil {
		return err
	}
	logger.Info("User authenticated")
	a.Scs.Put(c.Ctx(), string(ctxKeyAuthenticatedUserID), id)
	return c.Redirect("/")
}

func (a *API) Logout(c *dove.Context) error {
	if err := a.Scs.RenewToken(c.Ctx()); err != nil {
		return err
	}
	a.Scs.Remove(c.Ctx(), string(ctxKeyAuthenticatedUserID))
	return c.Redirect("/login")
}

func State(a *API, r *http.Request) *views.State {
	return &views.State{
		Version:       a.Version,
		Authenticated: a.IsAuthenticated(r),
	}
}

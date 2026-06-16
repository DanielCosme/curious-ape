package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"danicos.dev/daniel/curious-ape/assets"
	"danicos.dev/daniel/curious-ape/pkg/application"
	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/dove"
	"danicos.dev/daniel/curious-ape/pkg/oak"
	"danicos.dev/daniel/curious-ape/pkg/ui"
)

func Routes(a *API) http.Handler {
	d := dove.New(a.App.Log.Handler())

	// e.StaticFS("/static", echo.MustSubFS(views.StaticFS, "static"))

	d.Use(dove.MiddlewarePanicRecover)

	d.Prefix("/assets").GET(a.ServeStaticAssets)

	if a.App.Env == application.Dev {
		d.Use(DevMiddleware)
	}

	d.Use(a.MiddlewareLoadCookie)
	d.Use(a.MiddlewareAuthenticateFromSession)

	d.Endpoint("/login").
		GET(a.GetLoginForm).
		POST(a.Login).
		DELETE(a.Logout)

	d.Endpoint("/api/oauth2/fitbit/success").GET(a.FitbitSuccess)
	d.Endpoint("/api/oauth2/google/success").GET(a.GoogleSuccess)

	d.Use(a.MiddlewareRequireAuthentication)

	d.Endpoint("/").GET(a.Home)
	d.Endpoint("/habit/flip").PUT(a.HabitFlip)
	d.Endpoint("/day/sync").POST(a.DaySync)
	d.Endpoint("/integration").GET(a.IntegrationGet)
	d.Endpoint("/integrations").GET(a.IntegrationsGetList)
	d.Endpoint("/habits").GET(a.Habits)
	d.Endpoint("/sleep").GET(a.Sleep)
	d.Endpoint("/fitness").GET(a.Fitness)
	d.Endpoint("/deep_work").GET(a.DeepWork)
	d.Endpoint("/deadlines").GET(a.DeadlinesList)
	d.Endpoint("/deadline").
		GET(a.DeadlinesGetForm).
		POST(a.DeadlinesPostForm)

	return d
}

func (a *API) DeadlinesList(c *dove.Context) error {
	res, err := a.App.DeadlineList()
	if err == nil {
		state := State(a, c.Req)
		state.Deadlines.DS = res
		return c.RenderOK(ui.Deadlines(state))
	}
	return err
}

func (a *API) DeadlinesGetForm(c *dove.Context) error {
	state := State(a, c.Req)
	return c.RenderOK(ui.DeadlineForm(state))
}

func (a *API) DeadlinesPostForm(c *dove.Context) error {
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

func (a *API) DeepWork(c *dove.Context) error {
	days, err := a.App.DaysMonth(c.Ctx(), getDateParam(c))
	if err == nil {
		state := State(a, c.Req)
		state.Days = days
		return c.RenderOK(ui.DeepWork(state))
	}
	return err
}

func (a *API) Fitness(c *dove.Context) error {
	days, err := a.App.DaysMonth(c.Ctx(), getDateParam(c))
	if err == nil {
		state := State(a, c.Req)
		state.Days = days
		return c.RenderOK(ui.Fitness(state))
	}
	return err
}

func (a *API) Habits(c *dove.Context) error {
	state := State(a, c.Req)
	today := core.NewDate(time.Now())
	for _, month := range today.Months() {
		t := time.Date(today.Time().Year(), month+1, -1, 0, 0, 0, 0, time.UTC)
		d := core.NewDate(t)
		if month == today.Time().Month() {
			d = today
		}
		days, err := a.App.DaysMonthASC(c.Ctx(), d)
		if err == nil {
			state.DaysYear = append(state.DaysYear, days)
		} else {
			return err
		}
	}
	return c.RenderOK(ui.Habits(state))
}

func (a *API) Sleep(c *dove.Context) error {
	days, err := a.App.DaysMonth(c.Ctx(), getDateParam(c))
	if err == nil {
		state := State(a, c.Req)
		state.Days = days
		return c.RenderOK(ui.Sleep(state))
	}
	return err
}

func (a *API) ServeStaticAssets(c *dove.Context) (err error) {
	path, found := strings.CutPrefix(c.Req.URL.Path, "/assets/")
	if found {
		var data []byte
		if a.App.Env == application.Dev {
			// In dev: no-cache with revalidation so changes are picked up on refresh without hard-reload.
			c.Res.Header().Set("Cache-Control", "no-cache, must-revalidate")
			data, err = os.ReadFile("./assets/" + path)
		} else {
			// In prod: long-lived immutable cache (1 year). New deploys will use new asset content (ETag changes).
			c.Res.Header().Set("Cache-Control", "public, max-age=86400, immutable")
			data, err = assets.Assets.ReadFile(path)
		}
		if err == nil {
			mimeType := mime.TypeByExtension(filepath.Ext(path))
			c.Res.Header().Set("Content-Type", mimeType)
			hash := sha256.Sum256(data)
			etag := `"` + hex.EncodeToString(hash[:]) + `"` // Compute a strong ETag from content hash (works for both dev disk and prod embed)
			c.Res.Header().Set("ETag", etag)
			c.Res.Header().Set("X-Content-Type-Options", "nosniff")

			// Support conditional GET (304 Not Modified)
			if match := c.Req.Header.Get("If-None-Match"); match != "" && match == etag {
				c.Res.WriteHeader(http.StatusNotModified)
				return
			}
			_, err = c.Res.Write(data)
		}
		return
	} else {
		c.Res.WriteHeader(http.StatusNotFound)
		return errors.New(c.Req.URL.Path + " " + "not found")
	}
}

func (a *API) FitbitSuccess(c *dove.Context) error {
	c.ParseForm()
	return a.App.Oauth2Success(core.IntegrationFitbit, c.Req.FormValue("code"))
}

func (a *API) GoogleSuccess(c *dove.Context) error {
	c.ParseForm()
	return a.App.Oauth2Success(core.IntegrationGoogle, c.Req.FormValue("code"))
}

func (a *API) IntegrationGet(c *dove.Context) error {
	c.ParseForm()
	provider := c.Req.Form.Get("name")
	integrationInfo, err := a.App.IntegrationGet(c.Ctx(), core.Integration(provider))
	if err == nil {
		return c.RenderOK(ui.Integration(integrationInfo))
	}
	return err
}

func (a *API) IntegrationsGetList(c *dove.Context) error {
	integrationInfo, err := a.App.IntegrationsGetList()
	if err == nil {
		state := State(a, c.Req)
		state.Integrations = integrationInfo
		return c.RenderOK(ui.Integrations(state))
	}
	return err
}

func (a *API) Home(c *dove.Context) (err error) {
	days, err := a.App.DaysMonth(c.Ctx(), getDateParam(c))
	if err == nil {
		s := State(a, c.Req)
		s.Days = days
		return c.RenderOK(ui.Home(s))
	}
	return err
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

func (a *API) HabitFlip(c *dove.Context) error {
	c.ParseForm()
	id, err := strconv.Atoi(c.Req.Form.Get("id"))
	if err == nil {
		habit, err := a.App.HabitFlip(id)
		if err == nil {
			day, err := a.App.DayGetOrCreate(habit.Date)
			if err == nil {
				return c.RenderOK(ui.Day(day))
			}
		}
	}
	return err
}

func (a *API) DaySync(c *dove.Context) error {
	c.ParseForm()
	date, _ := core.NewDateFromISO8601(c.Req.Form.Get("date"))
	day, err := a.App.DaySync(c.Ctx(), date)
	if err == nil {
		return c.RenderOK(ui.Day(day))
	}
	return err
}

func (a *API) GetLoginForm(c *dove.Context) error {
	if a.IsAuthenticated(c.Req) {
		return c.Redirect("/")
	}
	return c.RenderOK(ui.Login(State(a, c.Req)))
}

func (a *API) Login(c *dove.Context) error {
	logger := oak.FromContext(c.Ctx())
	c.ParseForm()
	username := c.Req.PostFormValue("username")
	password := c.Req.PostFormValue("password")
	id, err := a.App.Authenticate(username, password)
	if err == nil {
		err = a.Scs.RenewToken(c.Ctx())
		if err == nil {
			logger.Info("User authenticated")
			a.Scs.Put(c.Ctx(), string(ctxKeyAuthenticatedUserID), id)
			return c.Redirect("/")
		}
	}
	return err
	// if errors.Is(err, persistence.ErrInvalidCredentials) {
	// 	// TODO: send http.StatusUnauthorized
	// 	return err
	// } else {
	// 	// TODO: send http.InternalServerError
	// 	return err
	// }
}

func (a *API) Logout(c *dove.Context) error {
	if err := a.Scs.RenewToken(c.Ctx()); err != nil {
		return err
	}
	a.Scs.Remove(c.Ctx(), string(ctxKeyAuthenticatedUserID))
	return c.Redirect("/login")
}

func State(a *API, r *http.Request) *ui.State {
	return &ui.State{
		Version:       a.Version,
		Authenticated: a.IsAuthenticated(r),
		CurrentPath:   r.URL.Path,
	}
}

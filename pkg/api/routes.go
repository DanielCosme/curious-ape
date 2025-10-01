package api

import (
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/dove"
)

func Routes(a *API) http.Handler {
	d := dove.New(a.App.Log.Handler())

	// e.StaticFS("/static", echo.MustSubFS(views.StaticFS, "static"))

	d.Use(dove.MiddlewarePanicRecover)

	d.Endpoint("/").GET(a.Home)

	// NOTE: Configure Logger, here?

	return d
}

func (a *API) Home(c *dove.Context) error {
	days, err := a.App.DaysMonth(c.Ctx(), core.NewDate(time.Now()))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, days)
}

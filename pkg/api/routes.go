package api

import (
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/pkg/core"
	"github.com/danielcosme/curious-ape/pkg/dove"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func Routes(a *API) http.Handler {
	d := dove.New()

	// e.StaticFS("/static", echo.MustSubFS(views.StaticFS, "static"))

	d.Use(dove.MiddlewarePanicRecover)
	d.Use(dove.MiddlewareLogRequest)

	d.Endpoint("/").GET(a.Home)

	return d
}

func (a *API) Home(c dove.Context) error {
	oak.Info("Home called")
	days, err := a.App.DaysMonth(core.NewDate(time.Now()))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, days)
}

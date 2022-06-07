package router

import (
	"net/http"

	"github.com/danielcosme/curious-ape/internal/api/router/middleware"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/rest"
)

func Routes(a *application.App) http.Handler {
	h := Handler{App: a}

	mux := http.NewServeMux()
	md := rest.NewMiddleware()

	md.Use(middleware.Logger(a))
	md.Use(middleware.RecoverPanic(a))
	md.Use(rest.MiddlewareParseForm)

	mux.HandleFunc("/", h.NotFound)
	mux.HandleFunc("/ping", h.Ping)

	mux.Handle("/habits", rest.NewMiddleware(middleware.SetDay(a), middleware.SetHabit(a)).Then(h.Habits))
	mux.HandleFunc("/habits/categories", h.HabitCategories)

	mux.HandleFunc("/days", h.DaysGetAll)

	mux.Handle("/sleep", rest.NewMiddleware(middleware.SetDay(a)).Then(h.SleepLogs))

	mux.HandleFunc("/oauth2/fitbit/connect", h.Oauth2FitbitConnect)
	mux.HandleFunc("/oauth2/fitbit/success", h.Oauth2FitbitSuccess)

	return md.Commit(mux)
}

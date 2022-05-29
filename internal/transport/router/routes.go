package router

import (
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/internal/transport/router/middleware"
	md "github.com/danielcosme/curious-ape/rest/middleware"
	"net/http"
)

func Routes(a *application.App) http.Handler {
	h := Handler{App: a}

	mux := http.NewServeMux()
	m := md.New()

	m.Use(md.LogRequest)
	m.Use(md.RecoverPanic)
	m.Use(md.Misc)

	mux.HandleFunc("/", h.NotFound)
	mux.HandleFunc("/ping", h.Ping)

	mux.Handle("/habits", md.New(middleware.SetDay(a), middleware.SetHabit(a)).Then(h.Habits))
	mux.HandleFunc("/habits/categories", h.HabitCategories)

	mux.HandleFunc("/days", h.DaysGetAll)

	mux.HandleFunc("/sleep/debug", h.SleepDebug)

	mux.HandleFunc("/oauth2/fitbit/connect", h.Oauth2FitbitConnect)
	mux.HandleFunc("/oauth2/fitbit/success", h.Oauth2FitbitSuccess)

	return m.Commit(mux)
}

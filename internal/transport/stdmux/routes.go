package stdmux

import (
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/curious-ape/rest/middleware"
	"net/http"
)

func Routes(a *application.App) http.Handler {
	h := Handler{A: a}

	mux := http.NewServeMux()
	md := middleware.New()

	md.Use(middleware.LogRequest)
	md.Use(middleware.RecoverPanic)

	mux.HandleFunc("/ping", h.Ping)
	// mux.HandleFunc("/habits/", a.HabitsHandler)
	mux.HandleFunc("/", h.NotFound)

	return md.Commit(mux)
}

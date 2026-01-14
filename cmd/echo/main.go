package main

import (
	"net/http"

	"github.com/danielcosme/curious-ape/pkg/dove"
	"github.com/danielcosme/curious-ape/pkg/oak"
)

func main() {
	logger := oak.NewDefault()
	mux := dove.New(logger.Handler())
	mux.Endpoint("/").GET(func(c *dove.Context) error {
		return c.HTML([]byte(`<h2>Hello, from Kubernetes!</h2>`))
	})
	Addr := ":4001"
	s := http.Server{
		Addr:    ":4001",
		Handler: mux,
	}
	logger.Info("Listening on: " + Addr)
	if err := s.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}
}

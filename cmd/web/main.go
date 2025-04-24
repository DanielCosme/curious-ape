package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	port := "4000"

	// I need a http server that listens for requests, on the designated port.
	server := http.Server{
		Addr:    ":" + port,
		Handler: HTTPRouter{},
	}

	// TODO(daniel): handle graceful termination with server.Shutdown and server.Close.
	// 		These signals come from UNIX signals.
	slog.Info("server started", "port", port)
	server.ListenAndServe()
}

type HTTPRouter struct {
}

func (s HTTPRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(daniel): Implement router.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Hi"}`)
}

package transport

import (
	"fmt"
	"net/http"
)

func Router(t *Transport) http.Handler {
	router := http.NewServeMux()

	routes := Routes(t)
	for _, r := range routes {
		router.HandleFunc(r.pattern, r.handler)
	}
	router.HandleFunc("GET /api/routes", func(w http.ResponseWriter, r *http.Request) {
		var result string
		for _, route := range routes {
			result += fmt.Sprintf("%s\n", route.pattern)
		}
		result += "GET /api/routes"
		fmt.Fprintln(w, result)
	})

	return router
}

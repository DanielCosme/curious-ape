package middleware

import (
	"github.com/danielcosme/curious-ape/rest"
	"mime"
	"net/http"
)

func CheckJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get(rest.HeaderContentType)

		mt, _, err := mime.ParseMediaType(contentType)
		if err != nil || mt != "application/json" {
			rest.ErrResponse(w, r, http.StatusUnsupportedMediaType, "Content-Code header must be application/json")
			return
		}

		next.ServeHTTP(w, r)
	})
}

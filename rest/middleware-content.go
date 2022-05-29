package rest

import (
	"mime"
	"net/http"
)

func CheckJsonContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get(HeaderContentType)

		mt, _, err := mime.ParseMediaType(contentType)
		if err != nil || mt != "application/json" {
			ErrResponse(w, r, http.StatusUnsupportedMediaType, "content-type header must be application/json")
			return
		}

		next.ServeHTTP(w, r)
	})
}

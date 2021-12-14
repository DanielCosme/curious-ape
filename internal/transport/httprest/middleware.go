package httprest

import (
	"github.com/danielcosme/curious-ape/rest"
	"mime"
	"net/http"
)

func CheckJsonContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get(rest.HeaderContentType)

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				rest.JSON(w, http.StatusBadRequest, rest.Envelope{"error": "Malformed Content-Type header"})
				return
			}

			if mt != "application/json" {
				rest.JSON(w, http.StatusUnsupportedMediaType, rest.Envelope{"error": "Content-Type header must be application/json"})
				return
			}

			next.ServeHTTP(w,r)
		} else {
			rest.JSON(w, http.StatusUnsupportedMediaType, rest.Envelope{"error": "Content-Type header must be application/json"})
			return
		}
	})
}

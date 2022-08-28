package middleware

import (
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core/application"
	"github.com/danielcosme/go-sdk/errors"
	"github.com/danielcosme/rest"
	"net/http"
	"strconv"
	"time"
)

func Logger(a *application.App) rest.HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			t := time.Now()
			rw = rest.NewResponseWriterPlus(rw)
			properties := map[string]string{
				"Method": r.Method,
				"Path":   r.URL.Path,
			}

			next.ServeHTTP(rw, r)

			// After
			rwPlus, ok := rw.(*rest.ResponseWriterPlus)
			if !ok {
				a.Log.Fatal(errors.NewFatal("response writer plus assertion failed"))
				return
			}
			status := rwPlus.Status()

			properties["Status"] = fmt.Sprintf("[%s] %s", strconv.Itoa(status), http.StatusText(status))
			properties["Time"] = time.Now().Sub(t).String()
			if status < 400 {
				// from 200 to 400
				a.Log.InfoP("", properties)
			} else if status < 500 {
				msg := "no details available"
				// from 400 to 500
				if rwPlus.Err != nil {
					msg = rwPlus.Err.Error()
				}
				a.Log.WarningP(msg, properties)
			} else {
				// 500 and above
				a.Log.ErrorP(rwPlus.Err, properties)
			}
		})
	}
}

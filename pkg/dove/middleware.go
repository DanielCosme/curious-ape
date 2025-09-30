package dove

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/danielcosme/curious-ape/pkg/oak"
)

// MiddlewareFunc defines a function to process middleware.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func MiddlewarePanicRecover(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		defer func() {
			if r := recover(); r != nil {
				if r == http.ErrAbortHandler {
					panic(r)
				}
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 4096)
				length := runtime.Stack(stack, false)
				stack = stack[:length]

				msg := fmt.Sprintf("PANIC %v %s", err, stack)
				oak.Error(msg)
			}
		}()
		return next(c)
	}
}

func MiddlewareLogRequest(next HandlerFunc) HandlerFunc {
	return func(c Context) error {
		// TODO: Implement me
		next(c)
		oak.Info("Request Logging")
		return nil
	}
}

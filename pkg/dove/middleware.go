package dove

import (
	"fmt"
	"net/http"
	"runtime"
)

// MiddlewareFunc defines a function to process middleware.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

func MiddlewarePanicRecover(next HandlerFunc) HandlerFunc {
	return func(c *Context) error {
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
				c.Log.Error(msg)
				c.Res.WriteHeader(http.StatusInternalServerError)
			}
		}()
		return next(c)
	}
}

package dove

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/pkg/oak"
)

// TODO:
// - Add Middleware.
// 	- Rate limiter.

// Dove is the main framework Instance.
type Dove struct {
	logHandler slog.Handler // Log backend.
	logLevel   slog.Level
	middleware []MiddlewareFunc
	routes     map[string]*endpoint
}

func New(logHandler slog.Handler) *Dove {
	d := &Dove{
		// TODO: make log Level configurable
		logHandler: logHandler,
		logLevel:   oak.LevelTrace,
		routes:     map[string]*endpoint{},
	}
	return d
}

func (d *Dove) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	logger := oak.New(oak.NewQueuedHandler(d.logLevel)).Layer("web")
	c := NewContext(r, rw, logger)

	c.Log.Info(fmt.Sprintf("%s %s", c.Req.Method, c.Req.RequestURI))
	defer func() {
		if !c.Res.Commited {
			c.Res.WriteHeader(http.StatusOK)
		}

		if queue, ok := c.Log.Handler().(*oak.QueuedHandler); ok {
			queue.EndTrace()
			msg := fmt.Sprintf("%d %s", c.Res.StatusCode, http.StatusText(c.Res.StatusCode))
			duration := time.Since(c.StartTime).String()
			if c.Res.StatusCode < 400 {
				c.Log.Info(msg, "duration", duration)
			} else {
				c.Log.Error(msg, "duration", duration)
			}
			queue.Dequeue(d.logHandler)
		}
	}()

	endpoint, ok := d.routes[r.URL.Path]
	if !ok {
		c.Res.WriteHeader(http.StatusNotFound)
		return
	}

	handler, ok := endpoint.Handlers[r.Method]
	if !ok {
		c.Res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := handler(c)
	if err != nil {
		// TODO: Improve global error handler.
		c.Log.Error(err.Error())
		c.Res.WriteHeader(http.StatusInternalServerError)
	}
}

func (d *Dove) Endpoint(path string) *endpoint {
	e := Endpoint(path)
	e.middleware = d.middleware
	d.routes[path] = e
	return e
}

func (d *Dove) Use(ms ...MiddlewareFunc) {
	d.middleware = append(d.middleware, ms...)
}

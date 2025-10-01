package dove

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/pkg/oak"
)

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(c *Context) error

// Context represents state of the current HTTP request.
type Context struct {
	Req       *http.Request
	Res       *Response
	StartTime time.Time
	Log       *oak.Oak
}

func NewContext(req *http.Request, rw http.ResponseWriter, logger *oak.Oak) *Context {
	c := &Context{
		Req:       req,
		Res:       NewResponse(rw),
		StartTime: time.Now(),
		Log:       logger,
	}
	c.Req = c.Req.WithContext(oak.WithContext(c.Ctx(), c.Log))
	return c
}

func (c *Context) JSON(status int, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	c.Res.Header().Set("Content-Type", "application/json")
	_, err = c.Res.Write(bytes)
	return err
}

func (c *Context) Ctx() context.Context {
	return c.Req.Context()
}

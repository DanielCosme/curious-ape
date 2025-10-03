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
	c.Res.WriteHeader(status)
	c.Res.Header().Set("Content-Type", "application/json")
	_, err = c.Res.Write(bytes)
	return err
}

func (c *Context) HTML(body []byte) error {
	c.Res.Header().Set("Content-Type", "text/html")
	_, err := c.Res.Write(body)
	return err
}

func (c *Context) Render(status int, r Renderer) error {
	c.Res.Header().Set("Content-Type", "text/html")
	c.Res.WriteHeader(status)
	return r.Render(c.Res.Writer)
}

func (c *Context) RenderOK(r Renderer) error {
	return c.Render(http.StatusOK, r)
}

func (c *Context) Ctx() context.Context {
	return c.Req.Context()
}

func (c *Context) ParseForm() {
	err := c.Req.ParseForm()
	if err != nil {
		panic(err)
	}
}

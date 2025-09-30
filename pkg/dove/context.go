package dove

import (
	"encoding/json"
	"net/http"
)

// HandlerFunc defines a function to serve HTTP requests.
type HandlerFunc func(c Context) error

// Context represents state of the current HTTP request.
type Context struct {
	Req      *http.Request
	Res      http.ResponseWriter
	Response *Response
}

func (c *Context) JSON(status int, payload any) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = c.Res.Write(bytes)
	return nil
}

package fitbit

import (
	"net/http"
	"net/url"
)

type Client struct {
	*http.Client
}

func (c *Client) Call(method string, params url.Values, body, i interface{}) error {
	// make the body a reader
	// create headers
	// handle errors
	//
	url := ""
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	_, err = c.Do(req)
	return err
}

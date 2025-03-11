package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	*http.Client
}

func (c *Client) Call(method, path string, urlParams url.Values, i interface{}) error {
	reqURL := BaseURL + path
	if urlParams != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, urlParams.Encode())
	}

	// Make request
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if strconv.Itoa(res.StatusCode)[:1] != "2" {
		return c.catchGoogleError(body)
	}
	if !json.Valid(body) {
		// c.out.Write(body)
		return errors.New("response body is not valid json")
	}

	return json.Unmarshal(body, i)
}

func (c *Client) catchGoogleError(b []byte) error {
	// c.out.Write(b)
	return errors.New("google api error")
}

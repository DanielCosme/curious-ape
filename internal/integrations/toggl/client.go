package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseURL = "https://api.track.toggl.com"

type Client struct {
	*http.Client
	token string
	out   io.Writer
}

func (c *Client) Call(method, path string, urlParams url.Values, payload any) error {
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
	req.SetBasicAuth(c.token, "api_token")
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return c.catchTogglErr(body)

	}
	if !json.Valid(body) {
		c.out.Write(body)
		return errors.New("response body is not valid json")
	}

	return json.Unmarshal(body, payload)
}

func (c *Client) catchTogglErr(body []byte) error {
	c.out.Write(body)
	return nil
}

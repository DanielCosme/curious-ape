package fitbit

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

func (c *Client) Call(method string, path string, urlParams url.Values, i interface{}) error {
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
		return c.catchFitbitError(body)
	}
	if !json.Valid(body) {
		return errors.New("response body is not valid json")
	}

	return json.Unmarshal(body, i)
}

func (c *Client) catchFitbitError(b []byte) error {
	// TODO unmarshal json if any and return string error
	// c.out.Write(b)
	return errors.New("fitbit api error")
}

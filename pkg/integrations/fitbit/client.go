package fitbit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
)

type Client struct {
	*http.Client
}

func (c *Client) Call(method string, path string, urlParams url.Values, i any) error {
	reqURL := BaseURL + path
	if urlParams != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, urlParams.Encode())
	}
	req, err := http.NewRequest(method, reqURL, nil)
	if err == nil {
		req.Header.Set("accept", "application/json")
		res, err := c.Do(req)
		if err == nil {
			body, err := io.ReadAll(res.Body)
			if err == nil {
				if strconv.Itoa(res.StatusCode)[:1] == "2" {
					if json.Valid(body) {
						return json.Unmarshal(body, i)
					}
					return errors.New("response body is not valid json")
				}
				return c.catchFitbitError(body)
			}
		}
	}
	return err
}

func (c *Client) catchFitbitError(b []byte) error {
	slog.Error("Fitbit ERR", "message", string(b))
	return errors.New("fitbit api error check logs")
}

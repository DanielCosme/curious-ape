package hevy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"danicos.dev/daniel/curious-ape/pkg/oak"
)

const BaseURL = "https://api.hevyapp.com"

type Client struct {
	*http.Client
	apiKey string
}

func (c *Client) Call(method, path string, reqBody, resPayload any) error {
	var bd io.Reader
	if reqBody != nil {
		jsonBody, err := json.Marshal(reqBody)
		if err == nil {
			bd = bytes.NewReader(jsonBody)
		} else {
			return err
		}
	}

	reqURL := BaseURL + path
	// Make request
	req, err := http.NewRequest(method, reqURL, bd)
	if err == nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("api-key", c.apiKey)
		res, err := c.Do(req)
		if err == nil {
			body, err := io.ReadAll(res.Body)
			if err == nil {
				if res.StatusCode >= 200 && res.StatusCode < 300 {
					return json.Unmarshal(body, resPayload)
				}
				return c.catchHevyErr(body)
			}
		}
	}
	return err
}

func (c *Client) catchHevyErr(body []byte) error {
	oak.Error(string(body))
	return errors.New("Hevy api error")
}

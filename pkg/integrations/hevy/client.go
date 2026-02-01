package hevy

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/danielcosme/curious-ape/pkg/oak"
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
		if err != nil {
			return err
		}
		bd = bytes.NewReader(jsonBody)
	}

	reqURL := BaseURL + path
	// Make request
	req, err := http.NewRequest(method, reqURL, bd)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return c.catchHevyErr(body)
	}

	return json.Unmarshal(body, resPayload)
}

func (c *Client) catchHevyErr(body []byte) error {
	oak.Error(string(body))
	return errors.New("Hevy api error")
}

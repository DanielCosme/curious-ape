package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

const BaseURL = "https://api.track.toggl.com"

// in case of 4xx error - don't try another request with the same payload, inspect the response body, most of the time it has a readable message.
// in case of 5xx error - have a random delay before the next request.
// in case of 429 (Too Many Requests) - back off for a few minutes (you can expect a rate of 1req/sec to be available).
// in case of 410 (Gone) - don't try this endpoint again.
// in case of 402 (Payment required) - workspace should be upgraded to have access to said feature, don't repeat the request until that has happened.

type Client struct {
	*http.Client
	token       string
	workspaceID int
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
		slog.Error(string(body))
		return errors.New("toggl response body is not valid json")
	}

	return json.Unmarshal(body, resPayload)
}

func (c *Client) catchTogglErr(body []byte) error {
	slog.Error(string(body))
	return errors.New("toggl api error")
}

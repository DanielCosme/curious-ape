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

func (c *Client) Call(method string, resourceURI string, urlParams url.Values, i interface{}) error {
	reqURL := BaseURL + resourceURI
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

	// fmt.Println("Fitbit RES body", string(body))
	if strconv.Itoa(res.StatusCode)[:1] != "2" {
		return catchFitbitError(body)
	}
	if !json.Valid(body) {
		return errors.New("response body is not valid json")
	}

	return json.Unmarshal(body, i)
}

func catchFitbitError(b []byte) error {
	// TODO unmarshal json if any and return string error
	fmt.Println("FITBIT ERR", string(b))
	return errors.New("fitbit error")
}

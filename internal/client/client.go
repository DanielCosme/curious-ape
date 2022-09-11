package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/danielcosme/go-sdk/errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var DefaultService = &Service{
	Client: http.DefaultClient,
}

type Service struct {
	*http.Client
	BaseURL  string
	username string
	password string
}

func NewService() *Service {
	return &Service{Client: &http.Client{}}
}

func (s *Service) Auth(username, password string) *Service {
	s.username = username
	s.password = password
	return s
}

func (s *Service) Host(host string) *Service {
	s.BaseURL = host
	return s
}

func (s *Service) Call(method, path string, body, bind any, params url.Values) error {
	reqURL := s.BaseURL + path
	if params != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, params.Encode())
	}

	var reader io.Reader
	if body != nil {
		// marshal the body into a reader
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, reqURL, reader)
	if err != nil {
		return err
	}
	if s.username != "" && s.password != "" {
		req.SetBasicAuth(s.username, s.password)
	}

	req.Header.Set("accept", "application/json")
	res, err := s.Do(req)
	if err != nil {
		return err
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if strconv.Itoa(res.StatusCode)[:1] != "2" {
		if body == nil {
			return errors.New(res.Status)
		}
		return s.CatchErr(resBody)
	}
	if json.Valid(resBody) {
		return json.Unmarshal(resBody, bind)
	}

	return nil
}

func (s *Service) CatchErr(b []byte) error {
	e := fmt.Errorf("error: %s", string(b))
	return e
}

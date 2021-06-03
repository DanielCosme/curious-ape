package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/danielcosme/curious-ape/internal/auth"
	"github.com/danielcosme/curious-ape/internal/data"
)

var ErrNoRecord = errors.New("no record")

type FitnessProvider struct {
	Auth  *auth.AuthConfig
	Token *data.AuthTokenModel
	Scope string
}

func (fit *FitnessProvider) LogsRange(start, end string) ([]map[string]string, error) {
	params := map[string]string{
		"startTime":    formatDate(start, 1),
		"endTime":      formatDate(end, 24),
		"activityType": "97",
	}
	url := BaseUrl + "?" + auth.UrlEncode(params)
	res, err := fit.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(res, &jsonResponse)
	if err != nil {
		return nil, err
	}
	if len(jsonResponse) == 0 {
		return nil, ErrNoRecord
	}

	result := []map[string]string{}
	arr := jsonResponse["session"].([]interface{})

	for _, session := range arr {
		m := make(map[string]string)
		s := session.(map[string]interface{})
		for k, v := range s {
			if k == "description" || k == "id" {
				continue
			}

			str, ok := v.(string)
			if !ok {
				tmp, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				m["packageName"] = tmp["packageName"].(string)
				continue
			}

			m[k] = str
		}
		result = append(result, m)
	}

	return result, nil
}

func formatDate(date string, offset int) string {
	d, _ := time.Parse("2006-01-02", date)
	dur := time.Duration(offset)
	dISO, _ := d.Add(time.Hour * dur).MarshalText()
	return string(dISO)
}

func (fit *FitnessProvider) DayLog(date string) (map[string]string, error) {
	params := map[string]string{
		"startTime":    formatDate(date, 1),
		"endTime":      formatDate(date, 24),
		"activityType": "97",
	}
	url := BaseUrl + "?" + auth.UrlEncode(params)
	res, err := fit.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var jsonResponse map[string]interface{}
	response := make(map[string]string)
	err = json.Unmarshal(res, &jsonResponse)
	if err != nil {
		return nil, err
	}
	if len(jsonResponse) == 0 {
		return nil, ErrNoRecord
	}

	arr := jsonResponse["session"].([]interface{})
	if len(arr) == 0 {
		return nil, ErrNoRecord
	}

	m := arr[0].(map[string]interface{})
	for k, v := range m {
		str, ok := v.(string)
		if !ok {
			tmp, ok := v.(map[string]interface{})
			if !ok {
				continue
			}
			response["packageName"] = tmp["packageName"].(string)
			continue
		}
		response[k] = str
	}

	return response, nil
}

func (fit *FitnessProvider) makeRequest(url string) (body []byte, err error) {
	isExpired := true
	times := 0

	for isExpired && times < 2 {
		token, err := fit.Token.Get("google")
		if err != nil {
			return nil, err
		}

		res, err := fit.Auth.MakeRequest(url, token.AccessToken)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		s := res.StatusCode
		if s == http.StatusForbidden || s == http.StatusUnauthorized {
			log.Println("Status code", res.StatusCode, "\nERR", string(body))
			times++
			err := fit.RefreshToken(token.RefreshToken)
			if err != nil {
				// TODO keep track of the days that have no logs.
				return body, err
			}
		} else if res.StatusCode == http.StatusOK {
			log.Println("Request Successfully received")
			isExpired = false
		} else {
			return nil, fmt.Errorf("Error Code: %v\n%s", res.StatusCode, string(body))
		}
	}

	return body, nil
}

// Source: https://developers.google.com/identity/protocols/oauth2/web-server#httprest_1
// -- Making Requests --
// GET /drive/v2/files HTTP/1.1
// Host: www.googleapis.com
// Authorization: Bearer access_token

// authorization request params - query string GET
//	redirect_uri=
//	prompt=consent ?
//	response_type=code
//	client_id=
//	scope=https://www.googleapis.com/auth/fitness.activity.read
//	access_type=offline

// token exchange params - url encoded body POST
//	code=
//	redirect_uri=
//	client_id=
//	client_secret=
//	scope=<empty>?
//	grant_type=authorization_code

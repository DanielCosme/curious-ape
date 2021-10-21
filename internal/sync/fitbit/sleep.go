package fitbit

import (
	"encoding/json"
	"fmt"
	"github.com/danielcosme/curious-ape/internal/core"
	"io"
	"log"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/auth"
)

const (
	ZeroDay = "2021-01-01"
	BaseUrl = "https://api.fitbit.com/1.2/user/-/"
)

var ErrNoRecord = fmt.Errorf("Error procesing the logs range result")

type SleepProvider struct {
	Auth  *auth.AuthConfig
	Token core.AuthTokenModel
	Scope string
}

func (sc *SleepProvider) LogsRange(start, end string) (map[string][]byte, error) {
	url := fmt.Sprintf("%s%s/date/%s/%s.json", BaseUrl, sc.Scope, start, end)
	result, err := sc.makeRequest(url)
	if err != nil {
		return nil, err
	}
	var jsonResponse map[string][]interface{}
	response := map[string][]byte{}

	err = json.Unmarshal(result, &jsonResponse)
	if err != nil {
		return nil, err
	}

	for _, v := range jsonResponse["sleep"] {
		blob, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrNoRecord
		}

		key := blob["dateOfSleep"].(string)
		jsonBlob, err := json.Marshal(blob)
		if err != nil {
			return nil, err
		}

		response[key] = jsonBlob
	}

	return response, nil
}

func (sc *SleepProvider) DayLog(date string) ([]byte, error) {
	url := fmt.Sprintf("%s%s/date/%s.json", BaseUrl, sc.Scope, date)

	result, err := sc.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(result, &jsonResponse)
	if err != nil {
		return nil, err
	}

	arr := jsonResponse["sleep"].([]interface{})

	if len(arr) == 0 {
		return nil, fmt.Errorf("no %s log found for %s", sc.Scope, date)
	}

	response, err := json.Marshal(arr[0])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (sc *SleepProvider) makeRequest(url string) (body []byte, err error) {
	isExpired := true
	times := 0

	for isExpired && times < 2 {
		token, err := sc.Token.Get(sc.Auth.Provider)
		if err != nil {
			return nil, err
		}

		res, err := sc.Auth.MakeRequest(url, token.AccessToken)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if res.StatusCode == http.StatusUnauthorized {
			log.Println("Status code", res.StatusCode, "\nERR", string(body))
			times++
			err := sc.refreshToken()
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

func (sc *SleepProvider) refreshToken() (err error) {
	log.Println("Refreshing Token")
	t, err := sc.Token.Get(sc.Auth.Provider)
	if err != nil {
		return err
	}

	newT, err := sc.Auth.RefreshToken(t.RefreshToken)
	if err != nil {
		return err
	}

	err = sc.Token.Update(newT)
	return nil
}

func (col *SleepProvider) AuthorizationURI() string {
	params := map[string]string{
		"client_id":     col.Auth.ClientID,
		"response_type": "code",
		"scope":         "sleep",
		"expires_in":    "604800",
		"redirect_uri":  col.Auth.RedirectURL,
	}
	urlEncoded := col.Auth.AuthorizationURL + "?" + auth.UrlEncode(params)
	return urlEncoded
}

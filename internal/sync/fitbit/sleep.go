package fitbit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/auth"
	"github.com/danielcosme/curious-ape/internal/data"
)

const (
	ZeroDay = "2020-02-01"
	BaseUrl = "https://api.fitbit.com/1.2/user/-/"
)

type SleepCollector struct {
	Auth  *auth.AuthConfig
	Token *data.AuthTokenModel
	Scope string
}

func (sc *SleepCollector) DayLog(date string) ([]byte, error) {
	url := fmt.Sprintf("%s%s/date/%s.json", BaseUrl, sc.Scope, date)

	log.Println("Getting daily Log")
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
	response, err := json.Marshal(arr[0])
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (sc *SleepCollector) makeRequest(url string) (body []byte, err error) {
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
			err = sc.refreshToken()
			if err != nil {
				return body, nil
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

func (sc *SleepCollector) refreshToken() (err error) {
	log.Println("Refreshing Token")
	t, err := sc.Token.Get(sc.Auth.Provider)
	if err != nil {
		return err
	}

	newT, err := sc.Auth.RefreshToken(t.RefreshToken)
	err = sc.Token.Update(newT)
	return err
}

func (col *SleepCollector) AuthorizationURI() string {
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

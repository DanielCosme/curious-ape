package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/danielcosme/curious-ape/internal/auth"
)

var GoogleAuth = &auth.AuthConfig{}

const BaseUrl = "https://www.googleapis.com/fitness/v1/users/me/sessions"

func init() {
	GoogleAuth.AuthorizationURL = "https://accounts.google.com/o/oauth2/v2/auth?"
	GoogleAuth.TokenRequestURL = "https://oauth2.googleapis.com/token"
	GoogleAuth.ClientSecret = os.Getenv("GOOGLE_SECRET")
	GoogleAuth.ClientID = os.Getenv("GOOGLE_ID")
	GoogleAuth.RedirectURL = os.Getenv("GOOGLE_REDIRECT_URI")
	GoogleAuth.Provider = "google"
}

func (fit *FitnessProvider) AuthorizationURI() string {
	params := map[string]string{
		"client_id":     fit.Auth.ClientID,
		"redirect_uri":  fit.Auth.RedirectURL,
		"response_type": "code",
		"scope":         "https://www.googleapis.com/auth/fitness.activity.read",
		"prompt":        "consent",
		"access_type":   "offline",
	}
	urlEncoded := fit.Auth.AuthorizationURL + auth.UrlEncode(params)
	return urlEncoded
}

func (fit *FitnessProvider) ExchangeCodeForToken(code string) (payload *auth.Token, err error) {
	payload, err = fit.tokens(code, "authorization")
	if err != nil {
		return nil, err
	}
	return payload, err
}

func (fit *FitnessProvider) RefreshToken(refreshToken string) (err error) {
	payload, err := fit.tokens(refreshToken, "refresh")
	if err != nil {
		return err
	}
	token, err := fit.Token.Get("google")
	if err != nil {
		return err
	}
	log.Println("cur token", token.AccessToken)
	log.Println("new token", payload.AccessToken)
	token.AccessToken = payload.AccessToken
	err = fit.Token.Update(*token)
	if err != nil {
		return err
	}
	return nil
}

func (fit *FitnessProvider) tokens(codeOrToken, grant string) (*auth.Token, error) {
	var token *auth.Token
	var params map[string]string
	if grant == "authorization" {
		params = map[string]string{
			"client_id":     fit.Auth.ClientID,
			"client_secret": fit.Auth.ClientSecret,
			"redirect_uri":  fit.Auth.RedirectURL,
			"grant_type":    "authorization_code",
			"code":          codeOrToken,
		}
	} else if grant == "refresh" {
		params = map[string]string{
			"client_id":     fit.Auth.ClientID,
			"client_secret": fit.Auth.ClientSecret,
			"grant_type":    "refresh_token",
			"refresh_token": codeOrToken,
		}
	} else {
		return nil, fmt.Errorf("Invalid grant.")
	}

	body := auth.UrlEncode(params)
	req, err := fit.tokensRequest(body)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &token)
	if err != nil {
		return nil, err
	}
	token.Service = fit.Auth.Provider

	if res.StatusCode == http.StatusBadRequest {
		return nil, errors.New("Failed to fetch or refresh tokens")
	}

	return token, nil
}

func (fit *FitnessProvider) tokensRequest(body string) (*http.Request, error) {
	req, err := http.NewRequest("POST", fit.Auth.TokenRequestURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}

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
//	grant_type= authorization_code
// -- RESPONSE __
// {
//   "access_token": "1/fFAGRNJru1FTz70BzhT3Zg",
//   "expires_in": 3920,
//   "token_type": "Bearer",
//   "scope": "https://www.googleapis.com/auth/drive.metadata.readonly",
//   "refresh_token": "1//xEoDL4iW3cxlI7yDbSRFYNG01kVKM2C-259HOF2aQbI"
// }

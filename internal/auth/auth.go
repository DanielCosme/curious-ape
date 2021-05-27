package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthConfig struct {
	AuthorizationURL string
	TokenRequestURL  string
	RedirectURL      string
	ClientID         string
	ClientSecret     string
	Provider         string
}

type Token struct {
	Service      string `json:"-"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (auth *AuthConfig) ExchangeCodeForToken(code string) (payload Token, err error) {
	payload, err = auth.tokens(code, "authorization")
	return payload, err
}

func (auth *AuthConfig) RefreshToken(refreshToken string) (payload Token, err error) {
	payload, err = auth.tokens(refreshToken, "refresh")
	return payload, err
}

func (auth *AuthConfig) tokens(codeOrToken, grant string) (Token, error) {
	var token Token
	var params map[string]string
	if grant == "authorization" {
		params = map[string]string{
			"client_id":    auth.ClientID,
			"grant_type":   "authorization_code",
			"code":         codeOrToken,
			"redirect_uri": auth.RedirectURL,
		}
	} else if grant == "refresh" {
		params = map[string]string{
			"grant_type":    "refresh_token",
			"refresh_token": codeOrToken,
		}
	} else {
		return token, fmt.Errorf("Invalid grant.")
	}

	body := UrlEncode(params)
	req, err := tokensRequest(body, auth)
	if err != nil {
		return token, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(resBody, &token)
	if err != nil {
		return token, err
	}
	token.Service = auth.Provider

	return token, nil
}

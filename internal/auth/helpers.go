package auth

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

func UrlEncode(values map[string]string) string {
	p := url.Values{}
	for key, value := range values {
		p.Add(key, value)
	}
	return p.Encode()
}

func tokensRequest(body string, auth *AuthConfig) (*http.Request, error) {
	req, err := http.NewRequest("POST", auth.TokenRequestURL, strings.NewReader(body))
	e := encodeCredentials(auth.ClientID, auth.ClientSecret)
	req.Header.Add("Authorization", "Basic "+e)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}

func encodeCredentials(id, secret string) string {
	msg := []byte(id + ":" + secret)
	strEncoded := base64.StdEncoding.EncodeToString(msg)
	return strEncoded
}

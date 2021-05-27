package auth

import "net/url"

type AuthConfig struct {
	AuthorizationURL string
	TokenRequestURL  string
	RedirectURL      string
	ClientID         string
	ClientSecret     string
}

func UrlEncode(values map[string]string) string {
	p := url.Values{}
	for key, value := range values {
		p.Add(key, value)
	}
	return p.Encode()
}

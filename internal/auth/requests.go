package auth

import (
	"log"
	"net/http"
)

func (a *AuthConfig) MakeRequest(url, accessToken string) (*http.Response, error) {
	req, err := a.createRequest(url, accessToken)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("AUTH Making request")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (ac *AuthConfig) createRequest(url, accessToken string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	return req, err
}

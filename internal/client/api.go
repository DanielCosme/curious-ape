package client

import "net/http"

type Config struct {
	Credentials Credentials `json:"credentials"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApeAPI struct {
	C *Service
}

func Ping() error {
	return DefaultService.Call(http.MethodGet, "/ping", nil, nil, nil)
}

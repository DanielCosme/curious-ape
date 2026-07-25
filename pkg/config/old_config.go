package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"danicos.dev/daniel/curious-ape/pkg/oak"
	"golang.org/x/oauth2"
)

type OldConfig struct {
	Port     int `json:"port"`
	Database struct {
		DSN string `json:"dsn"`
	} `json:"database"`
	Integrations struct {
		Fitbit *Oauth2Config `json:"fitbit"`
		Google *Oauth2Config `json:"google"`
		Toggl  struct {
			Token       string `json:"api_token"`
			WorkspaceID int    `json:"workspace_id"`
		} `json:"toggl"`
		Hevy struct {
			ApiKey string `json:"api_key"`
		} `json:"hevy"`
	} `json:"integrations"`
	Environment Environment
	Admin       User `json:"admin"`
	User        User `json:"user"`
	Guest       User `json:"guest"`
}

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func ParseEnvironment(s string) (Environment, error) {
	switch Environment(s) {
	case Prod:
		return Prod, nil
	case Dev:
		return Dev, nil
	case CI:
	case "":
		e := errors.New("empty environment field")
		slog.Error(e.Error())
		return "", e
	}
	e := errors.New("ivalid environment value: " + s)
	slog.Error(e.Error())
	return "", e
}

type Oauth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	TokenURL     string   `json:"token_url"`
	AuthURL      string   `json:"auth_url"`
	AuthStyle    int      `json:"auth_style"`
	Scopes       []string `json:"scopes"`
}

func (o Oauth2Config) ToConf() *oauth2.Config {
	oak.Info("Loading Oauth2 configuration", "redirect", o.RedirectURL)
	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   o.AuthURL,
			TokenURL:  o.TokenURL,
			AuthStyle: oauth2.AuthStyle(o.AuthStyle), // Zero value means auto-detect.
		},
		RedirectURL: o.RedirectURL,
		Scopes:      o.Scopes,
	}
}

func ReadConfiguration(cfg *OldConfig) *OldConfig {
	var err error
	var rawFile []byte

	env, err := ParseEnvironment(os.Getenv(ENVIRONMENT))
	if err != nil {
		logFatal(fmt.Errorf("environment variable %s is empty", ENVIRONMENT))
	}
	cfg.Environment = env
	configPath := "config.json"
	rawFile, err = os.ReadFile(configPath)
	exitIfErr(err)
	oak.Info("Configuration file loaded", "path", configPath)

	err = json.Unmarshal(rawFile, cfg)
	exitIfErr(err)
	return cfg
}

func exitIfErr(err error) {
	if err != nil {
		logFatal(err)
	}
}

func logFatal(err error) {
	oak.Fatal("Fatal failure", "err", err.Error(), "stack", string(debug.Stack()))
	os.Exit(1)
}

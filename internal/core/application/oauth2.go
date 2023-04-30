package application

import (
	"context"
	"net/http"
	"strconv"

	"github.com/danielcosme/curious-ape/internal/core/database"
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"github.com/danielcosme/go-sdk/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

func (a *App) AuthenticateMe(password string) (int, error) {
	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{entity.ProviderSelf}})
	if err != nil {
		return 0, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(o.AccessToken), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, database.ErrInvalidCredentials
		}
		return 0, err
	}

	return o.ID, nil
}

func (a *App) SetPassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{entity.ProviderSelf}})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	if o == nil {
		return a.db.Oauths.Create(&entity.Oauth2{
			Provider:    entity.ProviderSelf,
			AccessToken: string(hash),
		})
	}

	o.AccessToken = string(hash)
	_, err = a.db.Oauths.Update(o)
	return err
}

func (a *App) Oauth2ConnectProvider(provider string) (string, error) {
	p := entity.IntegrationProvider(provider)
	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{p}})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return "", err
	}
	if o == nil {
		if err = a.db.Oauths.Create(&entity.Oauth2{Provider: p}); err != nil {
			return "", err
		}
	}

	config := a.oauth2GetConfigurationForProvider(p)

	opts := []oauth2.AuthCodeOption{}
	switch p {
	case entity.ProviderGoogle:
		opts = append(opts,
			oauth2.SetAuthURLParam("access_type", "offline"),
			oauth2.SetAuthURLParam("approval_prompt", "force"),
		)
	}

	uri := config.AuthCodeURL("", opts...)
	return uri, nil
}

func (a *App) Oauth2Success(provider, code string) error {
	p := entity.IntegrationProvider(provider)
	config := a.oauth2GetConfigurationForProvider(p)

	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{p}})
	if err != nil {
		return err
	}

	o.AccessToken = t.AccessToken
	o.RefreshToken = t.RefreshToken
	o.Expiration = t.Expiry
	o.TokenType = t.Type()

	_, err = a.db.Oauths.Update(o)
	return err
}

func (a *App) Oauth2GetClient(provider entity.IntegrationProvider) (*http.Client, error) {
	config := a.oauth2GetConfigurationForProvider(provider)

	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{provider}})
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{
		AccessToken:  o.AccessToken,
		RefreshToken: o.RefreshToken,
		Expiry:       o.Expiration,
		TokenType:    o.TokenType,
	}

	switch provider {
	case entity.ProviderFitbit:
		newToken, err := config.TokenSource(context.Background(), token).Token()
		// Check if token is still valid, if not refresh it
		if newToken.AccessToken != token.AccessToken {
			// If token was refreshed we persist the new token info
			_, err = a.db.Oauths.Update(&entity.Oauth2{
				ID:           o.ID,
				AccessToken:  newToken.AccessToken,
				RefreshToken: newToken.RefreshToken,
				TokenType:    newToken.TokenType,
				Expiration:   newToken.Expiry,
			})
			token = newToken
		}
		if err != nil {
			return nil, err
		}
	default:
		// by default, we assume these tokens do not expire
	}

	return config.Client(context.Background(), token), err
}

func (a *App) oauth2GetConfigurationForProvider(provider entity.IntegrationProvider) *oauth2.Config {
	var config *entity.Oauth2Config
	switch provider {
	case entity.ProviderFitbit:
		config = a.cfg.Fitbit
	case entity.ProviderGoogle:
		config = a.cfg.Google
	default:
		return nil
	}

	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
		RedirectURL: config.RedirectURL,
		Scopes:      config.Scopes,
	}
}

func (a *App) Oauth2AddAPIToken(token, provider string) (string, error) {
	if token == "" {
		return "", errors.New("token is empty")
	}

	switch entity.IntegrationProvider(provider) {
	case entity.ProviderToggl:
		o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return "", err
		}

		if o == nil {
			o = &entity.Oauth2{
				Provider:    entity.ProviderToggl,
				AccessToken: token,
				TokenType:   "Bearer",
			}
			if err := a.db.Oauths.Create(o); err != nil {
				return "", err
			}
		} else {
			o.AccessToken = token
			if _, err := a.db.Oauths.Update(o); err != nil {
				return "", err
			}
		}
		api := a.sync.TogglClient(o.AccessToken)
		me, err := api.Me.GetProfile()
		if err != nil {
			return "", err
		}

		ws, err := api.Workspace.Get()
		if err != nil {
			a.Log.Error(err)
		}
		if len(ws) != 1 {
			return "", errors.New("we only support one workspace for toggl")
		}

		w := ws[0]
		o.ToogglOrganizationID = strconv.Itoa(w.OrganizationID)
		o.ToogglWorkSpaceID = strconv.Itoa(w.ID)
		if _, err := a.db.Oauths.Update(o); err != nil {
			return "", err
		}

		return me.Fullname, nil
	default:
		return "", errors.New("invalid provider")
	}
}

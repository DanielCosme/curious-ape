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

func (a *App) SetPassword(name, password string, role entity.Role) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if name == "" {
		return errors.New("name cannot be empty")
	}

	u, err := a.db.Users.Get(entity.UserFilter{Role: role})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	if u == nil {
		return a.db.Users.Create(&entity.User{
			Name:     name,
			Password: string(hash),
			Role:     role,
		})
	}
	_, err = a.db.Users.Update(&entity.User{
		Name:     name,
		Password: string(hash),
		Role:     role,
	})
	return err
}

// Authenticate returns userID if succesfully authenticated.
func (a *App) Authenticate(username, password string) (int, error) {
	u, err := a.db.Users.Get(entity.UserFilter{Name: username})
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return 0, database.ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, database.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return u.ID, nil
}

func (a *App) Exists(id int) (bool, error) {
	_, err := a.db.Users.Get(entity.UserFilter{ID: id})
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *App) Oauth2ConnectProvider(provider string) (string, error) {
	p := entity.IntegrationProvider(provider)
	o, err := a.db.Auths.Get(entity.AuthFilter{Provider: []entity.IntegrationProvider{p}})
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return "", err
	}
	if o == nil {
		if err = a.db.Auths.Create(&entity.Auth{Provider: p}); err != nil {
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
	if config == nil {
		a.Log.Error(errors.New("oauth2 success: invalid provider: " + provider))
		return errors.New("something went wrong")
	}

	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	o, err := a.db.Auths.Get(entity.AuthFilter{Provider: []entity.IntegrationProvider{p}})
	if err != nil {
		return err
	}

	o.AccessToken = t.AccessToken
	o.RefreshToken = t.RefreshToken
	o.Expiration = t.Expiry
	o.TokenType = t.Type()

	_, err = a.db.Auths.Update(o)
	return err
}

func (a *App) Oauth2GetClient(provider entity.IntegrationProvider) (*http.Client, error) {
	config := a.oauth2GetConfigurationForProvider(provider)

	o, err := a.db.Auths.Get(entity.AuthFilter{Provider: []entity.IntegrationProvider{provider}})
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
			_, err = a.db.Auths.Update(&entity.Auth{
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

func (a *App) AuthAddAPIToken(token, provider string) (string, error) {
	if token == "" {
		return "", errors.New("token is empty")
	}

	switch entity.IntegrationProvider(provider) {
	case entity.ProviderToggl:
		o, err := a.db.Auths.Get(entity.AuthFilter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return "", err
		}

		if o == nil {
			o = &entity.Auth{
				Provider:    entity.ProviderToggl,
				AccessToken: token,
				TokenType:   "Bearer",
			}
			if err := a.db.Auths.Create(o); err != nil {
				return "", err
			}
		} else {
			o.AccessToken = token
			if _, err := a.db.Auths.Update(o); err != nil {
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
		if _, err := a.db.Auths.Update(o); err != nil {
			return "", err
		}

		return me.Fullname, nil
	default:
		return "", errors.New("invalid provider")
	}
}

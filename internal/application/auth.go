package application

import (
	"context"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/danielcosme/curious-ape/internal/core"
	"github.com/danielcosme/curious-ape/internal/database/gen/models"
	"golang.org/x/oauth2"
)

// Oauth2ConnectProvider creates the models.Auth resource in the database.
// Then loads the provider configuration to generate the AuthCodeURL
func (a *App) Oauth2ConnectProvider(provider core.Integration) (string, error) {
	if _, err := a.db.Auths.Upsert(&models.AuthSetter{
		Provider:    omit.From(string(provider)),
		AccessToken: omit.From(""),
	}); err != nil {
		return "", err
	}

	// TODO: Consider loading this information when starting the Server/Application and having it in memory.
	config := oauth2.Config{
		ClientID:     a.cfg.Fitbit.ClientID,
		ClientSecret: a.cfg.Fitbit.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   a.cfg.Fitbit.AuthURL,
			TokenURL:  a.cfg.Fitbit.TokenURL,
			AuthStyle: 0,
		},
		RedirectURL: a.cfg.Fitbit.RedirectURL,
		Scopes:      a.cfg.Fitbit.Scopes,
	}

	var opts []oauth2.AuthCodeOption
	switch provider {
	case core.IntegrationGoogle:
		opts = append(opts,
			oauth2.SetAuthURLParam("access_type", "offline"),
			oauth2.SetAuthURLParam("approval_prompt", "force"),
		)
	}

	// Then generate Auth code URI
	return config.AuthCodeURL("", opts...), nil
}

func (a *App) Oauth2Success(provider, code string) error {
	// TODO: Consider loading this information when starting the Server/Application and having it in memory.
	config := oauth2.Config{
		ClientID:     a.cfg.Fitbit.ClientID,
		ClientSecret: a.cfg.Fitbit.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   a.cfg.Fitbit.AuthURL,
			TokenURL:  a.cfg.Fitbit.TokenURL,
			AuthStyle: 0,
		},
		RedirectURL: a.cfg.Fitbit.RedirectURL,
		Scopes:      a.cfg.Fitbit.Scopes,
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	_, err = a.db.Auths.Upsert(&models.AuthSetter{
		Provider:     omit.From(provider),
		AccessToken:  omit.From(token.AccessToken),
		RefreshToken: omitnull.From(token.RefreshToken),
		TokenType:    omitnull.From(token.Type()),
		Expiration:   omitnull.From(token.Expiry),
	})
	a.Log.Info("Exchange successful", "token_type", token.TokenType)
	return err
}

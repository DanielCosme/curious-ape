package integrations

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielcosme/curious-ape/internal/core"

	"golang.org/x/oauth2"
)

type Integrations struct {
	fitbit *oauth2.Config
	google *oauth2.Config
}

func New(fitbit *oauth2.Config) *Integrations {
	return &Integrations{
		fitbit: fitbit,
	}
}

func (i *Integrations) GenerateOauth2URI(provider core.Integration) string {
	var opts []oauth2.AuthCodeOption
	var config *oauth2.Config
	switch provider {
	case core.IntegrationFitbit:
		config = i.fitbit
	case core.IntegrationGoogle:
		config = i.fitbit
		opts = append(opts,
			oauth2.SetAuthURLParam("access_type", "offline"),
			oauth2.SetAuthURLParam("approval_prompt", "force"),
		)
	default:
		return ""
	}
	return config.AuthCodeURL("", opts...)
}

func (i *Integrations) GetHttpClient(provider core.Integration, currentToken *oauth2.Token, updateFunc func(integration core.Integration, t *oauth2.Token) error) (res *http.Client, err error) {
	var newToken *oauth2.Token
	var config *oauth2.Config
	switch provider {
	case core.IntegrationFitbit:
		config = i.fitbit
	case core.IntegrationGoogle:
		config = i.google
	default:
		panic("not implemented: " + provider)
	}
	// Refresh token if necessary.
	newToken, err = config.TokenSource(context.Background(), currentToken).Token()
	if err != nil {
		return
	}
	if newToken.AccessToken != currentToken.AccessToken {
		// Update token in database.
		err = updateFunc(provider, newToken)
		if err != nil {
			return
		}
		currentToken = newToken
	}
	res = config.Client(context.Background(), currentToken)
	return
}

func (i *Integrations) ExchangeToken(provider core.Integration, code string) (res *oauth2.Token, err error) {
	var config *oauth2.Config
	switch provider {
	case core.IntegrationFitbit:
		config = i.fitbit
	case core.IntegrationGoogle:
		config = i.google
	default:
		return res, errors.New("non-implemented provider: " + string(provider))
	}
	return config.Exchange(context.Background(), code)
}

package integrations

import (
	"context"
	"errors"
	"github.com/danielcosme/curious-ape/pkg/integrations/toggl"
	"log/slog"
	"net/http"

	"github.com/danielcosme/curious-ape/pkg/core"

	"golang.org/x/oauth2"
)

type Integrations struct {
	TogglAPI *toggl.API
	fitbit   *oauth2.Config
	google   *oauth2.Config
	list     []core.Integration
}

func New(togglWorkspaceID int, togglToken string, fitbit, google *oauth2.Config) *Integrations {
	i := &Integrations{
		TogglAPI: toggl.NewApi(togglWorkspaceID, togglToken),
		fitbit:   fitbit,
		google:   google,
	}
	if i.fitbit != nil {
		i.list = append(i.list, core.IntegrationFitbit)
	}
	if i.google != nil {
		i.list = append(i.list, core.IntegrationGoogle)
	}
	if togglToken != "" {
		i.list = append(i.list, core.IntegrationToggl)
	}
	return i
}

func (i *Integrations) IntegrationsList() []core.Integration {
	return i.list
}

func (i *Integrations) GenerateOauth2URI(provider core.Integration) string {
	var opts []oauth2.AuthCodeOption
	var config *oauth2.Config
	switch provider {
	case core.IntegrationFitbit:
		config = i.fitbit
	case core.IntegrationGoogle:
		config = i.google
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
	var config *oauth2.Config
	switch provider {
	case core.IntegrationFitbit:
		config = i.fitbit
	case core.IntegrationGoogle:
		config = i.google
	default:
		panic("not implemented: " + provider)
	}
	if !currentToken.Valid() {
		slog.Info("Refreshing token", "provider", string(provider))
		// Refresh token.
		currentToken, err = config.TokenSource(context.Background(), currentToken).Token()
		if err != nil {
			return
		}
		// Update token in database.
		err = updateFunc(provider, currentToken)
		if err != nil {
			return
		}
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

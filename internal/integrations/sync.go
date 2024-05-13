package integrations

import (
	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
	"github.com/danielcosme/curious-ape/internal/integrations/google"
	"github.com/danielcosme/curious-ape/internal/integrations/toggl"
	"net/http"
)

type Sync struct {
	Fitbit *fitbit.API
	Google *google.API
	Toggl  *toggl.API
}

func NewSync() *Sync {
	return &Sync{}
}

func (s *Sync) FitbitClient(c *http.Client) *fitbit.API {
	return fitbit.NewAPI(c)
}

func (s *Sync) GoogleClient(c *http.Client) *google.API {
	return google.NewAPI(c)
}

func (s *Sync) TogglClient(token string) *toggl.API {
	return toggl.NewApi(token)
}

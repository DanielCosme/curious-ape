package integrations

import (
	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
	"github.com/danielcosme/curious-ape/internal/integrations/google"
	"github.com/danielcosme/curious-ape/internal/integrations/toggl"
	"github.com/danielcosme/go-sdk/log"
	"net/http"
)

type Sync struct {
	Fitbit *fitbit.API
	Google *google.API
	Toggl  *toggl.API
	logger *log.Logger
}

func NewSync(l *log.Logger) *Sync {
	return &Sync{logger: l}
}

func (s *Sync) FitbitClient(c *http.Client) *fitbit.API {
	return fitbit.NewAPI(c, s.logger)
}

func (s *Sync) GoogleClient(c *http.Client) *google.API {
	return google.NewAPI(c, s.logger)
}

func (s *Sync) TogglClient(token string) *toggl.API {
	return toggl.NewApi(token, s.logger)
}

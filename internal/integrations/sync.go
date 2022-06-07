package integrations

import (
	"github.com/danielcosme/curious-ape/internal/integrations/fitbit"
	"github.com/danielcosme/curious-ape/sdk/log"
	"net/http"
)

type Sync struct {
	Fitbit *fitbit.API
	logger *log.Logger
}

func NewSync(l *log.Logger) *Sync {
	return &Sync{logger: l}
}

func (s *Sync) FitbitClient(c *http.Client) *fitbit.API {
	s.Fitbit = fitbit.NewAPI(c, s.logger)
	return s.Fitbit
}

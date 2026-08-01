package integration

import (
	"fmt"
	"log/slog"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/integrations"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db   bob.DB
	nats *nats.Conn
	sync *integrations.Integrations
}

func NewService(sync *integrations.Integrations, db bob.DB, nc *nats.Conn) *Service {
	s := &Service{
		db:   db,
		nats: nc,
		sync: sync,
	}
	if sync == nil {
		panic("sync cannot be nil")
	}

	return s
}

func (s *Service) GetIntegration(provider string) (res core.IntegrationInfo, err error) {
	var info []string
	// var authURL string
	// today := core.NewDateToday()
	status := core.IntegrationStatusNotImplemented

	switch provider {
	case core.IntegrationHevy:
		count, err := s.sync.Hevy.Workouts.Count()
		if err != nil {
			status = core.IntegrationStatusDisconnected
			info = append(info, err.Error())
		} else {
			status = core.IntegrationStatusConnected
			info = append(info, fmt.Sprintf("Number of workouts: %d", count))
		}
		res = core.IntegrationInfo{
			Name:   "Hevy",
			Info:   info,
			Status: status,
		}
	case core.IntegrationGoogle:
		// _, err := a.fitnessLogsFromGoogle(today)
		// if err != nil {
		// 	authURL = a.sync.GenerateOauth2URI(provider)
		// 	if authURL != "" {
		// 		status = core.IntegrationStatusDisconnected
		// 	}
		// 	info = append(info, err.Error())
		// } else {
		// 	status = core.IntegrationStatusConnected
		// }
		// res = core.IntegrationInfo{
		// 	Name:    "Google",
		// 	Info:    info,
		// 	AuthURL: authURL,
		// 	Status:  status,
		// }
		slog.Info("Google not implemented")
	case core.IntegrationFitbit:
		// sls, err := a.sleepLogsGetFromFitbit(today)
		// if err != nil {
		// 	authURL = a.sync.GenerateOauth2URI(provider)
		// 	if authURL != "" {
		// 		status = core.IntegrationStatusDisconnected
		// 	}
		// 	info = append(info, err.Error())
		// } else {
		// 	status = core.IntegrationStatusConnected
		// 	if len(sls) > 0 {
		// 		info = append(info, fmt.Sprintf("Total time asleep last night: %s", sls[0].TimeAsleep.String()))
		// 	}
		// }
		// res = core.IntegrationInfo{
		// 	Name:    "Fitbit",
		// 	Info:    info,
		// 	AuthURL: authURL,
		// 	Status:  status,
		// }
		slog.Info("Fitbit not implemented")
	case core.IntegrationToggl:
		profile, err := s.sync.TogglAPI.Me.GetProfile()
		if err != nil {
			info = append(info, err.Error())
		} else if profile != nil {
			status = core.IntegrationStatusConnected
			name := fmt.Sprintf("Profile name: %s", profile.FullName)
			timeZone := fmt.Sprintf("Timezone: %s", profile.Timezone)
			info = append(info, name, timeZone)

			ws, err := s.sync.TogglAPI.Workspace.Get()
			if err == nil {
				for _, w := range ws {
					info = append(info, fmt.Sprintf("Workspace: %s - ID: %d", w.Name, w.ID))
				}
			} else {
				slog.Error(err.Error())
			}

			summary, err := s.sync.TogglAPI.Reports.GetDaySummary(time.Now())
			if err != nil {
				slog.Error(err.Error())
			} else {
				info = append(info, fmt.Sprintf("Total time worked so far: %s", summary.TotalDuration))
			}
		}
		res = core.IntegrationInfo{
			Name:   "Toggl",
			Info:   info,
			Status: status,
		}
	default:
	}
	slog.Debug("integration: "+res.Name, "status", res.Status)

	return res, nil
}

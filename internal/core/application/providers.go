package application

import (
	"github.com/danielcosme/curious-ape/internal/core/entity"
	"strings"
)

func (a *App) TogglGetProjects() ([]entity.Entity, error) {
	var es []entity.Entity
	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
	if err != nil {
		return es, err
	}

	togglAPI := a.Sync.TogglClient(o.AccessToken)
	ps, err := togglAPI.Projects.GetAll(o.ToogglWorkSpaceID)
	if err != nil {
		return es, err
	}
	for _, p := range ps {
		es = append(es, entity.Entity{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return es, nil
}

// TogglAssignProjectsToGoal . Projects as (string) id, id, id
func (a *App) TogglAssignProjectsToGoal(ids string) error {
	o, err := a.db.Oauths.Get(entity.Oauth2Filter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
	if err != nil {
		return err
	}

	api := a.Sync.TogglClient(o.AccessToken)

	projects := strings.Split(ids, ",")
	for _, id := range projects {
		p, err := api.Projects.GetByID(o.ToogglWorkSpaceID, id)
		if err != nil {
			return err
		}
		a.Log.Tracef("Project from toggl successfully linked, name: %s", p.Name)
	}
	a.Log.Tracef("All projects for the goal exist on Toggl")

	// Update the project IDs
	o.ToogglProjectIDs = ids
	_, err = a.db.Oauths.Update(o)
	return err
}

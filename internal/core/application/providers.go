package application

import (
	"strings"

	"github.com/danielcosme/curious-ape/internal/core/entity"
)

func (a *App) TogglGetProjects() ([]entity.Entity, error) {
	var es []entity.Entity
	o, err := a.db.Auths.Get(entity.AuthFilter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
	if err != nil {
		return es, err
	}

	togglAPI := a.sync.TogglClient(o.AccessToken)
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
	o, err := a.db.Auths.Get(entity.AuthFilter{Provider: []entity.IntegrationProvider{entity.ProviderToggl}})
	if err != nil {
		return err
	}

	api := a.sync.TogglClient(o.AccessToken)

	projects := strings.Split(ids, ",")
	for _, id := range projects {
		p, err := api.Projects.GetByID(o.ToogglWorkSpaceID, id)
		if err != nil {
			return err
		}
		a.Log.TraceP("Project from time tracking provider successfully linked", props{
			"provider": entity.ProviderToggl,
			"name":     p.Name,
		})
	}
	a.Log.Trace("All projects for the goal exist on Toggl")

	// Update the project IDs
	o.ToogglProjectIDs = ids
	_, err = a.db.Auths.Update(o)
	return err
}

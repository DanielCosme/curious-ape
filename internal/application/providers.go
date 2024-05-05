package application

import (
	entity2 "github.com/danielcosme/curious-ape/internal/entity"
	"strings"
)

func (a *App) TogglGetProjects() ([]entity2.Entity, error) {
	var es []entity2.Entity
	o, err := a.db.Auths.Get(entity2.AuthFilter{Provider: []entity2.IntegrationProvider{entity2.ProviderToggl}})
	if err != nil {
		return es, err
	}

	togglAPI := a.sync.TogglClient(o.AccessToken)
	ps, err := togglAPI.Projects.GetAll(o.ToogglWorkSpaceID)
	if err != nil {
		return es, err
	}
	for _, p := range ps {
		es = append(es, entity2.Entity{
			ID:   p.ID,
			Name: p.Name,
		})
	}

	return es, nil
}

// TogglAssignProjectsToGoal . Projects as (string) id, id, id
func (a *App) TogglAssignProjectsToGoal(ids string) error {
	o, err := a.db.Auths.Get(entity2.AuthFilter{Provider: []entity2.IntegrationProvider{entity2.ProviderToggl}})
	if err != nil {
		return err
	}

	api := a.sync.TogglClient(o.AccessToken)

	projects := strings.Split(ids, ",")
	for _, id := range projects {
		_, err := api.Projects.GetByID(o.ToogglWorkSpaceID, id)
		if err != nil {
			return err
		}
	}

	// Update the project IDs
	o.ToogglProjectIDs = ids
	_, err = a.db.Auths.Update(o)
	return err
}

package toggl

import (
	"fmt"
	"net/http"
	"time"
)

type Project struct {
	ID          int       `json:"id"`
	WorkspaceID int       `json:"workspace_id"`
	Name        string    `json:"name"`
	IsPrivate   bool      `json:"is_private"`
	Active      bool      `json:"active"`
	At          time.Time `json:"at"`
	CreatedAt   time.Time `json:"created_at"`
	Color       string    `json:"color"`
	Recurring   bool      `json:"recurring"`
	ActualHours int       `json:"actual_hours"`
	Wid         int       `json:"wid"`
}

type ProjectsService struct {
	client *Client
}

func (s *ProjectsService) GetByID(workspaceID, projectID string) (*Project, error) {
	var project *Project
	path := fmt.Sprintf("/api/v9/workspaces/%s/projects/%s", workspaceID, projectID)
	err := s.client.Call(http.MethodGet, path, nil, &project)
	return project, err
}

func (s *ProjectsService) GetAll(workspaceID string) ([]*Project, error) {
	var projects []*Project
	path := fmt.Sprintf("/api/v9/workspaces/%s/projects", workspaceID)
	err := s.client.Call(http.MethodGet, path, nil, &projects)
	return projects, err
}

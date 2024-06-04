package google

import (
	"net/http"
)

const BaseURL = "https://www.googleapis.com"

type API struct {
	Fitness *FitnessService
}

func NewAPI(client *http.Client) API {
	c := API{Fitness: &FitnessService{client: Client{Client: client}}}
	return c
}

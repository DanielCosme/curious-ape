package google

import (
	"io"
	"net/http"
)

const BaseURL = "https://www.googleapis.com"

type API struct {
	Fitness *FitnessService
}

func NewAPI(client *http.Client, out io.Writer) *API {
	c := &API{Fitness: &FitnessService{client: Client{Client: client, out: out}}}
	return c
}

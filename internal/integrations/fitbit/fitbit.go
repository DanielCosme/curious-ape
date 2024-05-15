package fitbit

import (
	"net/http"
	"time"
)

const BaseURL = "https://api.fitbit.com"

// Fitbit-Rate-Limit-Limit: 	The quota number of calls.
// Fitbit-Rate-Limit-Remaining: The number of calls remaining before hitting the rate limit.
// Fitbit-Rate-Limit-Reset: 	The number of seconds until the rate limit resets.

type API struct {
	Sleep  *SleepService
	client *http.Client
}

func NewAPI(client *http.Client) API {
	return API{
		client: client,
		Sleep:  &SleepService{client: Client{Client: client}},
	}
}

func formatDate(time time.Time) string {
	return time.Format("2006-01-02")
}

package fitbit

import "net/http"

const BaseURL = "https://api.fitbit.com"
const DateFormat = "2006-01-02"

// Fitbit-Rate-Limit-Limit: The quota number of calls.
// Fitbit-Rate-Limit-Remaining: The number of calls remaining before hitting the rate limit.
// Fitbit-Rate-Limit-Reset: The number of seconds until the rate limit resets.

type API struct {
	Sleep *SleepService
}

func NewAPI(client *http.Client) *API {
	c := &API{Sleep: &SleepService{client: Client{client}}}
	return c
}

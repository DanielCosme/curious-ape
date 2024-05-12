package core

type DataSource string

const (
	WebUI  DataSource = "web_ui"
	Manual DataSource = "manual"
	Fitbit DataSource = "fitbit"
	Google DataSource = "google"
	Toggl  DataSource = "toggl"
)

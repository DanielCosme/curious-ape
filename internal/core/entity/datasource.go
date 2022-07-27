package entity

type DataSource string

const (
	Manual DataSource = "manual"
	Fitbit DataSource = "fitbit"
	Google DataSource = "google"
	Toggl  DataSource = "toggl"
)

func (ds DataSource) Str() string {
	return string(ds)
}

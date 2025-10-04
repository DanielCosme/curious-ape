package core

import "strings"

const (
	ISO8601           = "2006-01-02"
	HumanDate         = "02 Mon"
	HumanDateWeekDay  = "Monday, 02 Jan 2006"
	HumanDateWithTime = "Monday, 02 Jan 2006 at 15:04"
	Timestamp         = "15:04:05"
	Time              = "15:04"
)

type AuthRole string

const (
	AuthRoleAdmin AuthRole = "admin"
	AuthRoleUser  AuthRole = "user"
	AuthRoleGuest AuthRole = "guest"
)

type LogOrigin string

const (
	OriginLogManual LogOrigin = "manual"
	OriginLogWebUI  LogOrigin = "web_ui"
	OriginLogFitbit           = IntegrationFitbit
	OriginLogToggl            = IntegrationToggl
	OriginLogGoogle           = IntegrationGoogle
)

type Integration string

const (
	IntegrationFitbit = "fitbit"
	IntegrationGoogle = "google"
	IntegrationToggl  = "toggl"
)

func ToUpperFist(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

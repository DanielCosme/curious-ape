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

type HabitState string

const (
	HabitStateDone    HabitState = "done"
	HabitStateNotDone HabitState = "not_done"
	HabitStateNoInfo  HabitState = "no_info"
)

type HabitType string

const (
	HabitTypeEatHealthy HabitType = "food"
	HabitTypeWakeUp     HabitType = "wake_up"
	HabitTypeFitness    HabitType = "fitness"
	HabitTypeDeepWork   HabitType = "deep_work"
)

func HabitTypeFromString(s string) HabitType {
	switch s {
	case "food":
		return HabitTypeEatHealthy
	case "wake_up":
		return HabitTypeWakeUp
	case "fitness":
		return HabitTypeFitness
	case "deep_work":
		return HabitTypeDeepWork
	}
	return "undefined"
}

type AuthRole string

const (
	AuthRoleAdmin AuthRole = "admin"
	AuthRoleUser  AuthRole = "user"
	AuthRoleGuest AuthRole = "guest"
)

type OriginLog string

const (
	OriginLogManual OriginLog = "manual"
	OriginLogWebUI  OriginLog = "web_ui"
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

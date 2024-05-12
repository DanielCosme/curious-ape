package core

const (
	ISO8601           = "2006-01-02"
	HumanDate         = "02 Mon, Jan"
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
	HabitTypeExercise   HabitType = "fitness"
	HabitTypeDeepWork   HabitType = "deep_work"
	HabitTypeDynamic    HabitType = "dynamic"
)

type AuthRole string

const (
	AdminRole AuthRole = "admin"
	UserRole  AuthRole = "user"
	GuestRole AuthRole = "guest"
)
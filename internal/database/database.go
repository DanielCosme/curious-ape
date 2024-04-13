package database

type Repository struct {
	Days        Day
	Habits      Habit
	Auths       Authentication
	Users       User
	SleepLogs   SleepLog
	FitnessLogs FitnessLog
}

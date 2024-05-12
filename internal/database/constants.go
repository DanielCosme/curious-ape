package database

type Relation int

const (
	RelationDay Relation = iota + 1
	RelationHabit
	RelationHabitLogs
	RelationHabitCategory
	RelationSleep
	RelationFitness
)

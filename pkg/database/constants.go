package database

type Relation int

const (
	RelationDay Relation = iota + 1
	RelationHabit
	RelationSleep
	RelationFitness
	RelationWork
)

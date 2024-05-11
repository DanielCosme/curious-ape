package database

type Relation int

const (
	RelationHabit Relation = iota + 1
	RelationSleep
	RelationFitness
)

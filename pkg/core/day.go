package core

import (
	"bytes"
	"encoding/gob"
)

type DayRelations int

const (
	DayRelationHabits DayRelations = iota
	DayRelationFitnessLogs
	DayRelationDeepWorkLogs
	DayRelationSleepLogs
)

type Day struct {
	RepositoryCommon
	Date         Date
	Habits       DayHabits
	SleepLogs    []SleepLog
	FitnessLogs  []FitnessLog
	DeepWorkLogs []DeepWorkLog
}

type DayHabits struct {
	Score    int
	Hs       []Habit
	Sleep    Habit
	Fitness  Habit
	DeepWork Habit
	Eat      Habit
}

func (d *Day) IsZero() bool {
	return d.Date.Time().IsZero()
}

type DayParams struct {
	ID    uint
	Date  Date
	Dates DateSlice
	Order OrderParam
}

type DeadlineParams struct {
	Order OrderParam
}

type OrderParam int

const (
	ASC OrderParam = iota
	DESC
)

func DayRelationsAll() []DayRelations {
	return []DayRelations{
		DayRelationHabits,
		DayRelationFitnessLogs,
		DayRelationDeepWorkLogs,
		DayRelationSleepLogs,
	}
}

func (d *Day) Endocde() []byte {
	// NOTE: consider passing the encoder as parameter.
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(d); err != nil {
		panic("day: encoding failure: " + err.Error())
	}
	return buf.Bytes()
}

func Decode(data []byte) Day {
	var d Day
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&d); err != nil {
		panic("date: decode feilure: " + err.Error())
	}
	return d
}

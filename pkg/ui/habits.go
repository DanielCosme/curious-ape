package ui

import (
	"fmt"
	"git.danicos.dev/daniel/curious-ape/pkg/core"
	. "maragu.dev/gomponents"

	// ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

func Habits(s *State) Node {
	categories := []struct {
		name string
		typ  core.HabitType
	}{
		{"Wake Up", core.HabitTypeWakeUp},
		{"Fitness", core.HabitTypeFitness},
		{"Deep Work", core.HabitTypeDeepWork},
		{"Eat Healthy", core.HabitTypeEatHealthy},
	}

	daysCount := len(s.Days)
	gridStyle := fmt.Sprintf("grid-template-columns: 120px repeat(%d, 1fr);", daysCount)

	var nodes []Node
	nodes = append(nodes, Div(Class("grid-header"), Text("Category")))
	for _, day := range s.Days {
		nodes = append(nodes, Div(Class("grid-header"), Text(day.Date.Time().Format("02"))))
	}
	for _, cat := range categories {
		nodes = append(nodes, Div(Class("habit-category"), Text(cat.name)))
		for _, day := range s.Days {
			state := getHabitState(day, cat.typ)
			class := "habit-cell habit-" + string(state)
			nodes = append(nodes, Div(Class(class)))
		}
	}

	gridAttrs := []Node{Class("habits-grid"), Style(gridStyle)}
	allNodes := append(gridAttrs, nodes...)

	node := Div(
		H2(Text(s.Days[0].Date.Time().Month().String())),
		Div(allNodes...),
	)
	return layout("Habits", s, node)
}

func getHabitState(day core.Day, typ core.HabitType) core.HabitState {
	switch typ {
	case core.HabitTypeWakeUp:
		return day.Habits.Sleep.State
	case core.HabitTypeFitness:
		return day.Habits.Fitness.State
	case core.HabitTypeDeepWork:
		return day.Habits.DeepWork.State
	case core.HabitTypeEatHealthy:
		return day.Habits.Eat.State
	default:
		return core.HabitStateNoInfo
	}
}

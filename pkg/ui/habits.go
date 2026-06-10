package ui

import (
	"danicos.dev/daniel/curious-ape/pkg/core"
	"fmt"
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

	finalNodes := []Node{}
	for _, days := range s.DaysYear {
		var monthScore int
		daysCount := len(days)
		gridStyle := fmt.Sprintf("grid-template-columns: 120px repeat(%d, 1fr);", daysCount)

		var nodes []Node
		nodes = append(nodes, Div(Class("grid-header"), Text("Category")))
		for _, day := range days {
			monthScore += day.Habits.Score
			nodes = append(nodes, Div(Class("habit-grid-item grid-header"), Text(day.Date.Time().Format("02"))))
		}
		for _, cat := range categories {
			nodes = append(nodes, Div(Class("habit-category"), Text(cat.name)))
			for _, day := range days {
				state := getHabitState(day, cat.typ)
				class := "habit-grid-item habit-cell habit-" + string(state)
				nodes = append(nodes, Div(Class(class)))
			}
		}

		gridAttrs := []Node{Class("habits-grid"), Style(gridStyle)}
		allNodes := append(gridAttrs, nodes...)

		maxScore := daysCount * 4
		percentage := (float32(monthScore) * float32(100)) / float32(maxScore)
		node := Div(
			H2(Style("display: inline-block"), Text(days[0].Date.Time().Month().String())),
			Span(Class("month-score"), Text(fmt.Sprintf("%.0f%% %d/%d", percentage, monthScore, maxScore))),
			Div(allNodes...),
		)
		finalNodes = append(finalNodes, node)
	}
	return layout("Habits", s, Group(finalNodes))
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

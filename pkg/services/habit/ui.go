package habit

import (
	"fmt"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"
	lucide "github.com/eduardolat/gomponents-lucide"

	. "maragu.dev/gomponents"

	ds "maragu.dev/gomponents-datastar"
	. "maragu.dev/gomponents/html"
)

type HabitsPageState struct {
	Months []Month
}

type Month struct {
	m       time.Month
	days    []core.Date
	wakeUp  []core.Habit
	fitness []core.Habit
	work    []core.Habit
	eat     []core.Habit
	score   int
}

type DayHabit struct{}

func HabitsPage(s *ui.UIState, habitState *HabitsPageState) Node {
	s.Title = "Habits"

	finalNodes := []Node{}
	for _, month := range habitState.Months {
		finalNodes = append(finalNodes, HabitsGrid(month))
	}
	return ui.UILayout(s, Group(finalNodes))
}

func HabitsGrid(month Month) Node {
	daysCount := len(month.days)
	maxScore := daysCount * 4
	percentage := (float32(month.score) * float32(100)) / float32(maxScore)
	gridStyle := fmt.Sprintf("grid-template-columns: 120px repeat(%d, 1fr);", daysCount)

	monthDiv := Div(
		ID(fmt.Sprintf("%d-%d", month.days[0].Time.Year(), month.m)),
		H2(Class("mono"), Text(month.m.String())),
		Span(Class("month-score"), Text(fmt.Sprintf("%.0f%% %d/%d", percentage, month.score, maxScore))),
		Div(
			Class("habits-grid"), Style(gridStyle),
			Div(Class("grid-header"), Text("Category")),
			UI_habitHeaders(month.days),
			UI_HabitCells(month.wakeUp, core.HabitTypeWakeUp),
			UI_HabitCells(month.fitness, core.HabitTypeFitness),
			UI_HabitCells(month.work, core.HabitTypeDeepWork),
			UI_HabitCells(month.eat, core.HabitTypeEatHealthy),
		),
	)
	return monthDiv
}

func UI_habitHeaders(ds []core.Date) Node {
	return Map(ds, func(day core.Date) Node {
		fullDate := day.Time.Format("Mon 02")
		return Div(
			Class("habit-grid-item grid-header"),
			Text(day.Time.Format("02")),
			Span(Class("habit-header-tooltip"), Text(fullDate)),
		)
	})
}

func UI_HabitCells(hs []core.Habit, habitType core.HabitType) Node {
	var habitTypeName string
	switch habitType {
	case core.HabitTypeWakeUp:
		habitTypeName = "Wake up"
	case core.HabitTypeFitness:
		habitTypeName = "Fitness"
	case core.HabitTypeDeepWork:
		habitTypeName = "Deep Work"
	case core.HabitTypeEatHealthy:
		habitTypeName = "Eat Healthy"
	default:
		panic(fmt.Sprintf("unexpected core.HabitType: %#v", habitType))
	}
	return Group([]Node{
		Div(Class("habit-category"), Text(habitTypeName)),
		Map(hs, func(h core.Habit) Node {
			return ui_habitCell(h, habitType)
		}),
	})
}

func ui_habitCell(h core.Habit, habitType core.HabitType) Node {
	flipAction := fmt.Sprintf("@put('/habits/%d/flip')", h.ID)
	class := "habit-grid-item habit-cell habit-" + string(h.State)

	return Div(
		Class(class),
		ID(fmt.Sprintf("habit-%d", h.ID)),
		ds.On("click", flipAction),
		Span(Class("habit-cell-icon"), habitTypeIcon(habitType)),
	)
}

func habitTypeIcon(t core.HabitType) Node {
	switch t {
	case core.HabitTypeWakeUp:
		return lucide.Bed()
	case core.HabitTypeFitness:
		return lucide.Dumbbell()
	case core.HabitTypeDeepWork:
		return lucide.UserCog()
	case core.HabitTypeEatHealthy:
		return lucide.Beef()
	default:
		return nil
	}
}

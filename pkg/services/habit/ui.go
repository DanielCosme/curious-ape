package habit

import (
	"fmt"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/ui"

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
}

type DayHabit struct{}

func HabitsPage(s *ui.UIState, habitState *HabitsPageState) Node {
	s.Title = "Habits"

	finalNodes := []Node{}
	for _, month := range habitState.Months {
		daysCount := len(month.days) // Len of one of the Habits len...
		monthScore := 40             // TODO: Implement this!!
		maxScore := daysCount * 4
		percentage := (float32(monthScore) * float32(100)) / float32(maxScore)
		gridStyle := fmt.Sprintf("grid-template-columns: 120px repeat(%d, 1fr);", daysCount)

		monthDiv := Div(
			H2(Class("mono"), Text(month.m.String())),
			Span(Class("month-score"), Text(fmt.Sprintf("%.0f%% %d/%d", percentage, monthScore, maxScore))),
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
		finalNodes = append(finalNodes, monthDiv)
	}
	return ui.UILayout(s, Group(finalNodes))
}

func UI_habitHeaders(ds []core.Date) Node {
	return Map(ds, func(day core.Date) Node {
		return Div(Class("habit-grid-item grid-header"), Text(day.Time().Format("02")))
	})
}

func UI_HabitCells(hs []core.Habit, habitType core.HabitType) Node {
	return Group([]Node{
		Div(Class("habit-category"), Text(string(habitType))),
		Map(hs, func(h core.Habit) Node {
			return ui_habitCell(h)
		}),
	})
}

func ui_habitCell(h core.Habit) Node {
	flipAction := fmt.Sprintf("@put('/habits/%d/flip')", h.ID)
	class := "habit-grid-item habit-cell habit-" + string(h.State)

	return Div(
		Class(class),
		ID(fmt.Sprintf("habit-%d", h.ID)),
		ds.On("click", flipAction),
	)
}

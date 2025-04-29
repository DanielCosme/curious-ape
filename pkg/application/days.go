package application

import "time"

type Day struct {
	Date time.Time
	// Hardcoded Habits
	Wake    HabitDetail
	Workout HabitDetail
	Work    HabitDetail
	Eat     HabitDetail
}

type HabitDetail struct {
	State  HabitState
	Detail string
}

func (app *Application) Days() ([]Day, error) {
	return []Day{
		{
			Date:    time.Now(),
			Wake:    HabitDetail{State: HabitStateNotDone, Detail: "9:40"},
			Workout: HabitDetail{State: HabitStateDone, Detail: "13:30 - 14:30 (1h)"},
			Work:    HabitDetail{State: HabitStateDone, Detail: "5h45m"},
			Eat:     HabitDetail{State: HabitStateNoInfo},
		},
		{
			Date:    time.Now().AddDate(0, 0, 1),
			Wake:    HabitDetail{State: HabitStateDone, Detail: "5:45"},
			Workout: HabitDetail{State: HabitStateDone, Detail: "6:00 - 7:10 (1h10m)"},
			Work:    HabitDetail{State: HabitStateNotDone, Detail: "1h34m"},
			Eat:     HabitDetail{State: HabitStateNoInfo},
		},
	}, nil
}

func (app *Application) DaysGet(date time.Time) (*Day, error) {
	return &Day{
		Date:    date,
		Wake:    HabitDetail{State: HabitStateNotDone, Detail: "9:40"},
		Workout: HabitDetail{State: HabitStateDone, Detail: "13:30 - 14:30 (1h)"},
		Work:    HabitDetail{State: HabitStateDone, Detail: "5h45m"},
		Eat:     HabitDetail{State: HabitStateNoInfo},
	}, nil
}

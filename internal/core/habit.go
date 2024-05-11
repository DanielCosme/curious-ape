package core

type Habit struct {
	ID       int32
	Date     Date
	Category HabitCategory
	state    HabitState
	// Generated
	Day  *Day
	Logs []HabitLog
}

func NewHabit() *Habit {
	return &Habit{
		state: HabitStateNoInfo,
		Day:   nil,
		Logs:  nil,
	}
}

func (h Habit) State() HabitState {
	return h.state
}

func CalculateHabitStatus(logs []*HabitLog) HabitState {
	status := HabitStateNoInfo
	var override bool
	for _, log := range logs {
		if !log.IsAutomated {
			if log.Success {
				return HabitStateDone
			}
			status = HabitStateNotDone
			override = true
			continue
		}

		if !override {
			if log.Success {
				status = HabitStateDone
			} else if status != HabitStateDone {
				status = HabitStateNotDone
			}
		}
	}
	return status
}

type HabitCategory struct {
	ID          int32
	Name        string
	Type        HabitType
	Description string
}

type HabitLog struct {
	ID          int32
	Success     bool
	IsAutomated bool
	Origin      DataSource
	Note        string
}

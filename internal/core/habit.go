package core

type Habit struct {
	ID       int32
	DayID    int32
	Date     Date
	Category HabitCategory
	state    HabitState
	// Generated
	Day  *Day
	Logs []HabitLog
}

type HabitParams struct {
	Success    bool
	Date       Date
	CategoryID int32
	Origin     DataSource
	Automated  bool
}

func (hp *HabitParams) Valid() bool {
	return !hp.Date.Time().IsZero() && hp.CategoryID > 0
}

func NewHabit(d Date, c HabitCategory, logs []HabitLog) Habit {
	h := Habit{
		Date:     d,
		Category: c,
		Logs:     logs,
		state:    calculateHabitState(logs),
	}
	return h
}

func (h *Habit) State() HabitState {
	return h.state
}

func (h *Habit) IsZero() bool {
	return h.state == "" || h.Date.Time().IsZero()
}

func calculateHabitState(logs []HabitLog) HabitState {
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
}

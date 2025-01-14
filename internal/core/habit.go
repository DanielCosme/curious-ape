package core

type Habit struct {
	ID       int32
	DayID    int32
	Date     Date
	Category HabitCategory
	mainLog  HabitLog
	state    HabitState
	// Generated
	Day  *Day
	Logs []HabitLog
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
	Origin      OriginLog
	Detail      string // Wake Up
}

type NewHabitParams struct {
	Success   bool
	Date      Date
	HabitType HabitType
	Origin    OriginLog
	Automated bool
	Detail    string
}

func (hp *NewHabitParams) Valid() bool {
	return !hp.Date.Time().IsZero() && hp.HabitType != ""
}

func NewHabit(d Date, c HabitCategory, logs []HabitLog) Habit {
	h := Habit{
		Date:     d,
		Category: c,
		Logs:     logs,
	}
	h.state, h.mainLog = calculateHabitState(h.Logs)
	return h
}

func (h *Habit) State() HabitState {
	return h.state
}

func (h *Habit) Main() HabitLog {
	return h.mainLog
}

func (h *Habit) IsZero() bool {
	return h.state == "" || h.Date.Time().IsZero()
}

func calculateHabitState(logs []HabitLog) (state HabitState, mainLog HabitLog) {
	state = HabitStateNoInfo
	var override bool
	for _, log := range logs {
		if !log.IsAutomated {
			if log.Success {
				state = HabitStateDone
				mainLog = log
				return
			}
			state = HabitStateNotDone
			mainLog = log
			override = true
			continue
		}
		if !override {
			if log.Success {
				state = HabitStateDone
				mainLog = log
				return
			} else if state != HabitStateDone {
				mainLog = log
				state = HabitStateNotDone
			}
		}
	}
	return
}

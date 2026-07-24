package bobto

import (
	"log/slog"
	"time"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
)

func Day(d *models.Day) (day core.Day) {
	if d == nil {
		slog.Error("dayToCore: day is nil")
		return
	}
	day.ID = uint(d.ID)
	day.Date = core.NewDate(d.Date)
	for _, h := range d.R.Habits {
		habit := Habit(h)
		switch habit.Type {
		case core.HabitTypeWakeUp:
			day.Habits.Sleep = habit
		case core.HabitTypeFitness:
			day.Habits.Fitness = habit
		case core.HabitTypeDeepWork:
			day.Habits.DeepWork = habit
		case core.HabitTypeEatHealthy:
			day.Habits.Eat = habit
		}
		if h.State == string(core.HabitStateDone) {
			day.Habits.Score += 1
		}
		day.Habits.Hs = append(day.Habits.Hs, habit)
	}
	for _, sl := range d.R.SleepLogs {
		day.SleepLogs = append(day.SleepLogs, SleepLog(d, sl))
	}
	for _, fl := range d.R.FitnessLogs {
		day.FitnessLogs = append(day.FitnessLogs, FitnessLog(d, fl))
	}
	for _, wl := range d.R.DeepWorkLogs {
		day.DeepWorkLogs = append(day.DeepWorkLogs, WorkLog(d, wl))
	}
	return day
}

func Habit(h *models.Habit) (habit core.Habit) {
	if h == nil {
		slog.Error("habitToCore: habit is nil")
		return
	}
	habit.ID = uint(h.ID)
	habit.Date = core.NewDate(h.R.Day.Date)
	habit.State = core.HabitState(h.State)
	habit.Type = core.HabitType(h.R.HabitCategory.Kind)
	habit.Note = h.NOTE.GetOrZero()
	habit.Automated = h.Automated
	return
}

func FitnessLog(day *models.Day, bobFl *models.FitnessLog) (fl core.FitnessLog) {
	fl.ID = uint(bobFl.ID)
	fl.Date = core.NewDate(day.Date)
	fl.Title = bobFl.Title
	fl.StartTime = bobFl.StartTime
	fl.EndTime = bobFl.EndTime
	fl.Note = bobFl.Note
	fl.Type = core.TimelineTypeFitness
	fl.FitnessType = core.FitnessLogType(bobFl.Type)
	fl.Origin = core.LogOrigin(bobFl.Origin)
	return
}

func SleepLog(day *models.Day, s *models.SleepLog) (sl core.SleepLog) {
	sl = core.SleepLog{
		Date:        core.NewDate(day.Date),
		IsMainSleep: s.IsMainSleep.GetOrZero(),
		TimeAsleep:  time.Duration(s.TimeAsleep.GetOrZero()),
		TimeInBed:   time.Duration(s.TotalTimeInBed.GetOrZero()),
		TimelineLog: core.TimelineLog{
			Title:     s.Title,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
			Type:      core.TimelineTypeSleep,
			Note:      s.NOTE.GetOrZero(),
		},
	}
	sl.ID = uint(s.ID)
	return sl
}

func WorkLog(day *models.Day, bob *models.DeepWorkLog) (log core.DeepWorkLog) {
	log.ID = uint(bob.ID)
	log.Date = core.NewDate(day.Date)
	log.Title = bob.Title
	log.StartTime = bob.StartTime
	log.EndTime = bob.EndTime
	log.Note = bob.Note
	log.Type = core.TimelineTypeDeepWork
	return
}

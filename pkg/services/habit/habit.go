package habit

import (
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db bob.DB
}

func NewService(db bob.DB) *Service {
	return &Service{db: db}
}

func (s *Service) HabitUpsert(params core.Habit) (habit core.Habit, err error) {
	habit, err = upsert(s.db, params)
	if err == nil && habit.State != core.HabitStateNoInfo {
		slog.Info("Habit logged",
			"type", habit.Type,
			"state", habit.State,
			"date", habit.Date.Time().Format(core.HumanDateWeekDay),
		)
	}
	return
}

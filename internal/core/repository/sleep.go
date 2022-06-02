package repository

import "github.com/danielcosme/curious-ape/internal/core/entity"

type SleepLog interface {
	Create(*entity.SleepLog) error
	Update(*entity.SleepLog, ...entity.SleepLogJoin) (*entity.SleepLog, error)
	Get(entity.SleepLogFilter, ...entity.SleepLogJoin) (*entity.SleepLog, error)
	Find(entity.SleepLogFilter, ...entity.SleepLogJoin) ([]*entity.SleepLog, error)
	Delete(id int) error
	// Helpers
	ToIDs([]*entity.SleepLog) []int
	ToDayIDs([]*entity.SleepLog) []int
}

func ExecuteSleepLogPipeline(ssl []*entity.SleepLog, hjs ...entity.SleepLogJoin) error {
	if !(len(ssl) > 0) {
		return nil
	}

	for _, hj := range hjs {
		if err := hj(ssl); err != nil {
			return err
		}
	}
	return nil
}

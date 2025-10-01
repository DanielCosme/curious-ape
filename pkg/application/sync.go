package application

import (
	"github.com/danielcosme/curious-ape/pkg/core"
)

func (a *App) SyncDay(date core.Date) (core.Day, error) {
	errCh := make(chan error)

	go func() {
		errCh <- a.sleepSync(date)
	}()
	go func() {
		errCh <- a.fitnessSync(date)
	}()
	go func() {
		errCh <- a.deepWorkSync(date)
	}()
	for range 3 {
		if err := <-errCh; err != nil {
			a.Log.Error(err.Error())
		}
	}

	return a.DayGetOrCreate(date)
}

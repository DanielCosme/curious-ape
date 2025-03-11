package application

import (
	"github.com/danielcosme/curious-ape/pkg/core"
)

func (a *App) SyncDay(d core.Date) (core.Day, error) {
	errCh := make(chan error)

	go func() {
		errCh <- a.sleepSync(d)
	}()
	go func() {
		errCh <- a.fitnessSync(d)
	}()
	go func() {
		errCh <- a.deepWorkSync(d)
	}()
	for i := 0; i < 3; i++ {
		if err := <-errCh; err != nil {
			a.Log.Error(err.Error())
		}
	}

	return a.DayGetOrCreate(d)
}

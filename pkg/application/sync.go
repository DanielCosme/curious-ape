package application

import (
	"context"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"git.danicos.dev/daniel/curious-ape/pkg/oak"
)

func (a *App) DaySync(ctx context.Context, date core.Date) (core.Day, error) {
	logger := oak.FromContextWithLayer(ctx, "app")
	defer logger.PopLayer()

	errCh := make(chan error)

	go func() {
		errCh <- a.sleepSync(ctx, date)
	}()
	go func() {
		errCh <- a.fitnessSync(ctx, date)
	}()
	go func() {
		errCh <- a.deepWorkSync(ctx, date)
	}()
	for range 3 {
		if err := <-errCh; err != nil {
			logger.Error(err.Error())
		}
	}

	return a.DayGetOrCreate(date)
}

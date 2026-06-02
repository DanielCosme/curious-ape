package application

import (
	"context"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"git.danicos.dev/daniel/curious-ape/pkg/oak"
)

func (a *App) DeadlineCreate(ctx context.Context, params core.Deadline) (core.Deadline, error) {
	logger := oak.FromContext(ctx)

	err := params.Validate()
	if err != nil {
		return params, err
	}
	res, err := a.db.Deadlines.Create(params)
	if err != nil {
		return params, err
	}

	logger.Info("Deadline created",
		"Title", res.Title,
		"End Date", res.EndDate.String(),
		"recurring", res.Recurring,
	)
	return res, nil
}

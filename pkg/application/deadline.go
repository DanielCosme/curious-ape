package application

import (
	"context"
	"sort"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"git.danicos.dev/daniel/curious-ape/pkg/oak"
)

func (a *App) DeadlineCreate(ctx context.Context, params core.Deadline) (core.Deadline, error) {
	logger := oak.FromContext(ctx)
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

func (a *App) DeadlineList() ([]core.Deadline, error) {
	res, err := a.db.Deadlines.Find(core.DeadlineParams{})
	if err != nil {
		return nil, err
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].DaysLeft > res[j].DaysLeft // Change '>' to '<' for ascending
	})
	return res, nil
}

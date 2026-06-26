package application

import (
	"context"
	"sort"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/oak"
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

func (a *App) DeadlineList(ctx context.Context) ([]core.Deadline, error) {
	logger := oak.FromContext(ctx).Layer("app")
	defer logger.PopLayer()

	res, err := a.db.Deadlines.Find(core.DeadlineParams{})
	if err != nil {
		return nil, err
	}

	for idx, d := range res {
		if d.DaysLeft < 0 {
			if d.Recurring {
				// Add one year
				d.EndDate = core.NewDate(d.EndDate.Time().AddDate(1, 0, 0))
				err := a.db.Deadlines.Update(d)
				if err != nil {
					return nil, err
				}
				logger.Info("Recurring deadline updated",
					"title", d.Title,
					"End Date", d.EndDate.Time().Format(core.HumanDateWeekDay),
				)
				res[idx] = d
				continue
			}

			if err = a.db.Deadlines.Delete(d.ID); err != nil {
				return nil, err
			}
			logger.Info("Recurring deadline deleted", "title", d.Title)
			res[idx] = core.Deadline{}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].DaysLeft < res[j].DaysLeft // Change '>' to '<' for ascending
	})
	return res, nil
}

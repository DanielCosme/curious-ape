package persistence

import (
	"context"

	"git.danicos.dev/daniel/curious-ape/database/gen/models"
	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
)

type Deadlines struct {
	db bob.DB
}

func (d *Deadlines) Create(params core.Deadline) (deadlineRes core.Deadline, err error) {
	s := &models.DeadlineSetter{
		Title:     omit.From(params.Title),
		StartTime: omit.From(params.StartDate.Time()),
		EndTime:   omit.From(params.EndDate.Time()),
		Recurring: omit.From(params.Recurring),
	}
	deadline, err := models.Deadlines.Insert(s).One(context.Background(), d.db)
	if err != nil {
		return deadlineRes, catchDBErr("dealines: create", err)
	}
	return deadlineToCore(deadline), nil
}

func deadlineToCore(params *models.Deadline) core.Deadline {
	d := core.Deadline{
		Title:     params.Title,
		StartDate: core.NewDate(params.StartTime),
		EndDate:   core.NewDate(params.EndTime),
		Recurring: params.Recurring,
	}
	d.ID = uint(params.ID)
	return d
}

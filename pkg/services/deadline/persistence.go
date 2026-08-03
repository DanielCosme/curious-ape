package deadline

import (
	"context"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/gen/bob/models"
	"danicos.dev/daniel/curious-ape/pkg/persistence"
	"github.com/aarondl/opt/omit"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
)

func create(db bob.DB, params core.Deadline) (deadlineRes core.Deadline, err error) {
	err = params.Validate()
	if err != nil {
		return params, err
	}
	s := &models.DeadlineSetter{
		Title:     omit.From(params.Title),
		StartTime: omit.From(params.StartDate.Time),
		EndTime:   omit.From(params.EndDate.Time),
		Recurring: omit.From(params.Recurring),
	}
	deadline, err := models.Deadlines.Insert(s).One(context.Background(), db)
	if err != nil {
		return deadlineRes, persistence.CatchDBErr("dealines: create", err)
	}
	return deadlineToCore(deadline, core.NewDateToday()), nil
}

func find(db bob.DB, params core.DeadlineParams) (ds []core.Deadline, err error) {
	res, err := buildDeadlineQuery(params).All(context.Background(), db)
	if err != nil {
		return ds, persistence.CatchDBErr("deadlines: find", err)
	}
	today := core.NewDateToday()
	for _, deadline := range res {
		ds = append(ds, deadlineToCore(deadline, today))
	}
	return
}

func delete(db bob.DB, id uint) error {
	_, err := models.Deadlines.Delete(
		models.DeleteWhere.Deadlines.ID.EQ(int64(id)),
	).Exec(context.Background(), db)
	return persistence.CatchDBErr("deadlines: delete", err)
}

func update(db bob.DB, params core.Deadline) error {
	s := &models.DeadlineSetter{
		Title:     omit.From(params.Title),
		EndTime:   omit.From(params.EndDate.Time),
		Recurring: omit.From(params.Recurring),
	}
	_, err := models.Deadlines.Update(
		models.UpdateWhere.Deadlines.ID.EQ(int64(params.ID)),
		s.UpdateMod(),
	).Exec(context.Background(), db)
	return persistence.CatchDBErr("deadline: update", err)
}

func deadlineToCore(params *models.Deadline, today core.Date) core.Deadline {
	d := core.Deadline{
		Title:     params.Title,
		StartDate: core.NewDate(params.StartTime),
		EndDate:   core.NewDate(params.EndTime),
		Recurring: params.Recurring,
	}
	d.ID = uint(params.ID)
	d.DaysLeft = core.DaysLeft(today, d.EndDate)
	return d
}

func buildDeadlineQuery(p core.DeadlineParams) *sqlite.ViewQuery[*models.Deadline, models.DeadlineSlice] {
	q := models.Deadlines.Query()
	if p.Order == core.DESC {
		q.Apply(sm.OrderBy(models.Deadlines.Columns.StartTime).Desc())
	}
	return q
}

package deadline

import (
	"log/slog"
	"sort"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db   bob.DB
	nats *nats.Conn
}

func NewService(db bob.DB, nc *nats.Conn) *Service {
	s := &Service{
		db:   db,
		nats: nc,
	}
	return s
}

func (svc *Service) List() ([]core.Deadline, error) {
	res, err := find(svc.db, core.DeadlineParams{})
	if err != nil {
		return nil, err
	}

	for idx, d := range res {
		if d.DaysLeft < 0 {
			if d.Recurring {
				// Add one year
				d.EndDate = core.NewDate(d.EndDate.Time.AddDate(1, 0, 0))
				err := update(svc.db, d)
				if err != nil {
					return nil, err
				}
				slog.Info("Recurring deadline updated",
					"title", d.Title,
					"End Date", d.EndDate.Time.Format(core.HumanDateWeekDay),
				)
				res[idx] = d
				continue
			}

			if err = delete(svc.db, d.ID); err != nil {
				return nil, err
			}
			slog.Info("Recurring deadline deleted", "title", d.Title)
			res[idx] = core.Deadline{}
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].DaysLeft < res[j].DaysLeft // Change '>' to '<' for ascending
	})
	return res, nil
}

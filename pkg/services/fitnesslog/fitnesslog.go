package fitnesslog

import (
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db          bob.DB
	nats        *nats.Conn
	integration core.FitnessIntegration
}

func NewService(db bob.DB, nc *nats.Conn, integration core.FitnessIntegration) *Service {
	s := &Service{
		db:   db,
		nats: nc,
	}
	if nc != nil {
		nc.Subscribe(event.DaySync, s.listen)
	}
	return s
}

func (svc *Service) listen(msg *nats.Msg) {
	switch msg.Subject {
	case event.DaySync:
		err := svc.sync(core.DateDecode(msg.Data))
		if err != nil {
			slog.Error("Fitness Log: failed to synchronize", "err", err)
		}
	}
}

func (svc *Service) sync(date core.Date) error {
	logs, err := svc.integration.GetFitnessLogs(date)
	if err != nil {
		return err
	}

	// Upsert
	if len(logs) > 0 {
	}
	svc.nats.Publish(event.WorklogSynced, core.Encode(logs))
	return nil
}

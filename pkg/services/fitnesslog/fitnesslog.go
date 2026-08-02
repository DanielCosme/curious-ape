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
		db:          db,
		nats:        nc,
		integration: integration,
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
	result := core.LogSyncPayload{Date: date}
	logs, err := svc.integration.GetFitnessLogs(date)
	if err != nil {
		return err
	}

	for _, fl := range logs {
		_, err := upsert(svc.db, fl)
		if err != nil {
			return err
		}
		result.FitnessLogs = append(result.FitnessLogs, fl)
	}

	svc.nats.Publish(event.FitnesslogSynced, core.Encode(result))
	svc.nats.Publish(event.DaySynced, date.Enc())
	return nil
}

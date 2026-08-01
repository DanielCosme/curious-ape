package worklog

import (
	"fmt"
	"log/slog"

	"danicos.dev/daniel/curious-ape/pkg/core"
	"danicos.dev/daniel/curious-ape/pkg/event"
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db          bob.DB
	nats        *nats.Conn
	integration core.WorkLogIntegration
}

func NewService(db bob.DB, ns *nats.Conn, integration core.WorkLogIntegration) *Service {
	s := &Service{
		db:          db,
		nats:        ns,
		integration: integration,
	}
	if ns != nil {
		ns.Subscribe(event.DaySync, s.listen)
	}
	return s
}

func (svc *Service) listen(msg *nats.Msg) {
	switch msg.Subject {
	case event.DaySync:
		err := svc.sync(core.DateDecode(msg.Data))
		if err != nil {
			slog.Error("Work Log: failed to synchronize", "err", err)
		}
	}
}

func (svc *Service) sync(date core.Date) error {
	logs, err := svc.integration.GetDayEntries(date)
	if err != nil {
		return err
	}

	for _, log := range logs {
		_, err := upsert(svc.db, log)
		if err != nil {
			return err
		}
		duration := log.EndTime.Sub(log.StartTime)
		t := fmt.Sprintf("%s-%s (%s)", log.StartTime.Format(core.Time), log.EndTime.Format(core.Time), duration)
		slog.Info("Deep work log added: "+t, "origin", log.Origin)
	}
	if len(logs) > 0 {
		svc.nats.Publish(event.WorklogSynced, core.Encode(logs))
	}
	return nil
}

package sleeplog

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
	integration core.SleepIntegration
}

func NewService(db bob.DB, nc *nats.Conn, is core.SleepIntegration) *Service {
	s := &Service{
		db:          db,
		nats:        nc,
		integration: is,
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
			slog.Error("Sleep Log: failed to synchronize", "err", err)
		}
	}
}

func (svc *Service) sync(date core.Date) error {
	result := core.LogSyncPayload{Date: date}
	logs, err := svc.integration.GetSleepLogs(date)
	if err != nil {
		return err
	}
	for _, sl := range logs {
		_, err := upsert(svc.db, sl)
		if err != nil {
			return err
		}
		result.SleepLogs = append(result.SleepLogs, sl)
	}

	svc.nats.Publish(event.SleepLogSynced, core.Encode(result))
	svc.nats.Publish(event.DaySynced, date.Enc())
	return nil
}

package worklog

import (
	"github.com/nats-io/nats.go"
	"github.com/stephenafamo/bob"
)

type Service struct {
	db   bob.DB
	nats *nats.Conn
}

func NewService(db bob.DB, ns *nats.Conn) *Service {
	return &Service{db: db, nats: ns}
}

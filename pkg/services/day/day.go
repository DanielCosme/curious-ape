package day

import "github.com/stephenafamo/bob"

type Service struct {
	db bob.DB
}

func NewService(db bob.DB) *Service {
	return &Service{db: db}
}

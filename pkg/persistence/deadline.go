package persistence

import (
	"git.danicos.dev/daniel/curious-ape/pkg/core"
	"github.com/stephenafamo/bob"
)

type Deadlines struct {
	db bob.DB
}

func (d *Deadlines) Create(deadlineData core.Deadline) (deadlineRes core.Deadline, err error) {
	return
}

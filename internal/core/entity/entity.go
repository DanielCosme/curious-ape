package entity

import (
	"github.com/danielcosme/curious-ape/internal/utils/uuid"
	"time"
)

type Entity struct {
	ID        uint      `db:"id"`
	UUID      string    `db:"uuid"`
	CreatedAt time.Time `db:"creation_time"`
	UpdatedAt time.Time `db:"update_time"`
}

func NewEntity() Entity {
	return Entity{
		UUID:      uuid.NewUUID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

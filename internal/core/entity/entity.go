package entity

import (
	"time"

	"github.com/danielcosme/curious-ape/internal/utils/uuid"
)

type UUID string

type Entity struct {
	UUID
	ID        string
	NID       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func generateID() UUID {
	return UUID(uuid.NewUUID())
}

func NewEntity() *Entity {
	return &Entity{
		UUID:      generateID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

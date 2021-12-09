package entity

import (
	"time"

	"github.com/danielcosme/curious-ape/internal/utils/uuid"
)

type ID string

type Entity struct {
	ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func generateID() ID {
	return ID(uuid.NewUUID())
}

func NewEntity() *Entity {
	return &Entity{
		ID:        generateID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

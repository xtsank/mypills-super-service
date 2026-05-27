package dictionary

import (
	"github.com/google/uuid"
)

type Item struct {
	ID   uuid.UUID
	Name string
}

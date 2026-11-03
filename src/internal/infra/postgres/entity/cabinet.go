package entity

import "github.com/google/uuid"

type CabinetEntity struct {
	ID     uuid.UUID `db:"id"`
	UserID uuid.UUID `db:"user_id"`
	Name   string    `db:"name"`
}

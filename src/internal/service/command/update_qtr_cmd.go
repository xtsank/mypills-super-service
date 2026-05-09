package command

import (
	"github.com/google/uuid"
)

type UpdateQtyCmd struct {
	ID       uuid.UUID
	Quantity float32
}

func NewUpdateQtyCmd(id uuid.UUID, quantity float32) *UpdateQtyCmd {
	return &UpdateQtyCmd{
		ID:       id,
		Quantity: quantity,
	}
}

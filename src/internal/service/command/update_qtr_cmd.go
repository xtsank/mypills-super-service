package command

import (
	"github.com/google/uuid"
)

type UpdateQtyCmd struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Quantity float32
}

func NewUpdateQtyCmd(id uuid.UUID, quantity float32, userID ...uuid.UUID) *UpdateQtyCmd {
	cmd := &UpdateQtyCmd{
		ID:       id,
		Quantity: quantity,
	}
	if len(userID) > 0 {
		cmd.UserID = userID[0]
	}
	return cmd
}

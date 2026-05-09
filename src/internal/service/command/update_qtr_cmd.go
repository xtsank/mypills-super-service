package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

type UpdateQtyCmd struct {
	ID       uuid.UUID
	Quantity float32
}

func NewUpdateQtyCmd(id uuid.UUID, quantity float32) (*UpdateQtyCmd, error) {
	if quantity < 0 {
		return nil, errors.ErrQtyTooLow
	}

	return &UpdateQtyCmd{
		ID:       id,
		Quantity: quantity,
	}, nil
}

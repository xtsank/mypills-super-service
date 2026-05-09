package req

import "github.com/google/uuid"

type UpdateQtyDto struct {
	ID  uuid.UUID `json:"id" binding:"required"`
	Qty float32   `json:"qty" binding:"required"`
}

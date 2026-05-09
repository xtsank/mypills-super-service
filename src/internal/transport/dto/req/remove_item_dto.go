package req

import "github.com/google/uuid"

type RemoveItemDto struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

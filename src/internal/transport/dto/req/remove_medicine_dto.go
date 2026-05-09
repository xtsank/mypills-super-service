package req

import "github.com/google/uuid"

type RemoveMedicineDto struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

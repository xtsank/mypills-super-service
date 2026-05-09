package req

import "github.com/google/uuid"

type UpdateLinksDto struct {
	MedicineID uuid.UUID   `json:"medicine_id" binding:"required"`
	IDs        []uuid.UUID `json:"ids" binding:"required"`
}

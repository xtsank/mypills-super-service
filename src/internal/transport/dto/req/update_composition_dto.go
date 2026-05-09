package req

import "github.com/google/uuid"

type UpdateCompositionDto struct {
	MedicineID uuid.UUID             `json:"medicine_id" binding:"required"`
	Substances []*ActiveSubstanceDto `json:"substances" binding:"required"`
}

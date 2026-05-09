package req

import "github.com/google/uuid"

type SelectMedicineDto struct {
	IllnessID uuid.UUID `json:"illness_id" binding:"required"`
}

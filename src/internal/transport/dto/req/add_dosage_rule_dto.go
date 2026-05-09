package req

import "github.com/google/uuid"

type AddDosageRuleDto struct {
	MedicineID uuid.UUID      `json:"medicine_id" binding:"required"`
	Dosage     *DosageRuleDto `json:"dosage" binding:"required"`
}

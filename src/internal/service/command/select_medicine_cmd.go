package command

import "github.com/google/uuid"

type SelectMedicineCmd struct {
	UserID    uuid.UUID
	IllnessID uuid.UUID
}

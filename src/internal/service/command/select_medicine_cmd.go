package command

import "github.com/google/uuid"

type SelectMedicineCmd struct {
	UserID    uuid.UUID
	IllnessID uuid.UUID
}

func NewSelectMedicineCmd(userID, illnessID uuid.UUID) *SelectMedicineCmd {
	return &SelectMedicineCmd{
		UserID:    userID,
		IllnessID: illnessID,
	}
}

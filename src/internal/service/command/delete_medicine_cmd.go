package command

import "github.com/google/uuid"

type DeleteMedicineCmd struct {
	ID uuid.UUID
}

func NewDeleteMedicineCmd(id uuid.UUID) *DeleteMedicineCmd {
	return &DeleteMedicineCmd{ID: id}
}

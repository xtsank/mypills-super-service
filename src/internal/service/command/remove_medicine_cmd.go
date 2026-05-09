package command

import "github.com/google/uuid"

type RemoveMedicineCmd struct {
	ID uuid.UUID
}

func NewRemoveMedicineCmd(id uuid.UUID) *RemoveMedicineCmd {
	return &RemoveMedicineCmd{
		ID: id,
	}
}

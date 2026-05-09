package command

import "github.com/google/uuid"

type UpdateLinksCmd struct {
	MedicineID uuid.UUID
	IDs        []uuid.UUID
}

func NewUpdateLinksCmd(id uuid.UUID, ids []uuid.UUID) *UpdateLinksCmd {
	return &UpdateLinksCmd{
		MedicineID: id,
		IDs:        ids,
	}
}

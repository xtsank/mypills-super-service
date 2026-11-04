package command

import "github.com/google/uuid"

type DeleteCabinetCmd struct {
	UserID    uuid.UUID
	CabinetID uuid.UUID
}

func NewDeleteCabinetCmd(userID, cabinetID uuid.UUID) *DeleteCabinetCmd {
	return &DeleteCabinetCmd{UserID: userID, CabinetID: cabinetID}
}

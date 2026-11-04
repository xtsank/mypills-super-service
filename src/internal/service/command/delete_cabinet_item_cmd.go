package command

import "github.com/google/uuid"

type DeleteCabinetItemCmd struct {
	UserID    uuid.UUID
	CabinetID uuid.UUID
	ItemID    uuid.UUID
}

func NewDeleteCabinetItemCmd(userID, cabinetID, itemID uuid.UUID) *DeleteCabinetItemCmd {
	return &DeleteCabinetItemCmd{UserID: userID, CabinetID: cabinetID, ItemID: itemID}
}

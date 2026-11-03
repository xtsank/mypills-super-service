package command

import "github.com/google/uuid"

type UpdateRestockCmd struct {
	UserID          uuid.UUID
	ItemID          uuid.UUID
	DesiredQuantity float64
}

func NewUpdateRestockCmd(userID, itemID uuid.UUID, desiredQuantity float64) *UpdateRestockCmd {
	return &UpdateRestockCmd{UserID: userID, ItemID: itemID, DesiredQuantity: desiredQuantity}
}

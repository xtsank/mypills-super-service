package command

import "github.com/google/uuid"

type DeleteRestockCmd struct {
	UserID uuid.UUID
	ItemID uuid.UUID
}

func NewDeleteRestockCmd(userID, itemID uuid.UUID) *DeleteRestockCmd {
	return &DeleteRestockCmd{UserID: userID, ItemID: itemID}
}

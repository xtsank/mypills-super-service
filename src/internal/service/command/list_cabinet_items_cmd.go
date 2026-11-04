package command

import "github.com/google/uuid"

type ListCabinetItemsCmd struct {
	UserID            uuid.UUID
	CabinetID         uuid.UUID
	ID                *uuid.UUID
	LowStock          bool
	ExpiresWithinDays *int
	Expired           bool
	Limit             int
	Offset            int
}

func NewListCabinetItemsCmd(userID, cabinetID uuid.UUID, id *uuid.UUID, lowStock, expired bool,
	expiresWithinDays *int, limit, offset int) *ListCabinetItemsCmd {
	return &ListCabinetItemsCmd{UserID: userID, CabinetID: cabinetID, ID: id, LowStock: lowStock,
		Expired: expired, ExpiresWithinDays: expiresWithinDays, Limit: limit, Offset: offset}
}

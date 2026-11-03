package command

import "github.com/google/uuid"

type ListCabinetsCmd struct {
	UserID uuid.UUID
	ID     *uuid.UUID
	Limit  int
	Offset int
}

func NewListCabinetsCmd(userID uuid.UUID, id *uuid.UUID, limit, offset int) *ListCabinetsCmd {
	return &ListCabinetsCmd{UserID: userID, ID: id, Limit: limit, Offset: offset}
}

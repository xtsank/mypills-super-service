package command

import "github.com/google/uuid"

type ListRestockCmd struct {
	UserID uuid.UUID
	ID     *uuid.UUID
	Limit  int
	Offset int
}

func NewListRestockCmd(userID uuid.UUID, id *uuid.UUID, limit, offset int) *ListRestockCmd {
	return &ListRestockCmd{UserID: userID, ID: id, Limit: limit, Offset: offset}
}

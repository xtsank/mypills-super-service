package command

import "github.com/google/uuid"

type RemoveItemCmd struct {
	ID uuid.UUID
}

func NewRemoveItemCmd(id uuid.UUID) *RemoveItemCmd {
	return &RemoveItemCmd{ID: id}
}

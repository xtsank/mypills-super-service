package command

import "github.com/google/uuid"

type RemoveItemCmd struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func NewRemoveItemCmd(id uuid.UUID, userID ...uuid.UUID) *RemoveItemCmd {
	cmd := &RemoveItemCmd{ID: id}
	if len(userID) > 0 {
		cmd.UserID = userID[0]
	}
	return cmd
}

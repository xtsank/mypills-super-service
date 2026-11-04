package command

import "github.com/google/uuid"

type GetProfileCmd struct {
	UserID uuid.UUID
}

func NewGetProfileCmd(userID uuid.UUID) *GetProfileCmd {
	return &GetProfileCmd{UserID: userID}
}

package command

import "github.com/google/uuid"

type CreateCabinetCmd struct {
	UserID uuid.UUID
	Name   string
}

func NewCreateCabinetCmd(userID uuid.UUID, name string) *CreateCabinetCmd {
	return &CreateCabinetCmd{UserID: userID, Name: name}
}

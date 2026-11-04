package command

import "github.com/google/uuid"

type RenameCabinetCmd struct {
	UserID    uuid.UUID
	CabinetID uuid.UUID
	Name      string
}

func NewRenameCabinetCmd(userID, cabinetID uuid.UUID, name string) *RenameCabinetCmd {
	return &RenameCabinetCmd{UserID: userID, CabinetID: cabinetID, Name: name}
}

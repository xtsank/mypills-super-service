package command

import "github.com/google/uuid"

type CreateRestockCmd struct {
	UserID          uuid.UUID
	MedicineID      uuid.UUID
	DesiredQuantity float64
}

func NewCreateRestockCmd(userID, medicineID uuid.UUID, desiredQuantity float64) *CreateRestockCmd {
	return &CreateRestockCmd{UserID: userID, MedicineID: medicineID, DesiredQuantity: desiredQuantity}
}

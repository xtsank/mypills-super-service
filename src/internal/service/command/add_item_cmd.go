package command

import (
	"time"

	"github.com/google/uuid"
)

type AddItemCmd struct {
	UserID            uuid.UUID
	MedicineID        uuid.UUID
	DateOfManufacture time.Time
	Quantity          float32
}

func NewAddItemCmd(userId, medicineId uuid.UUID, date time.Time, quantity float32) *AddItemCmd {
	return &AddItemCmd{
		UserID:            userId,
		MedicineID:        medicineId,
		DateOfManufacture: date,
		Quantity:          quantity,
	}
}

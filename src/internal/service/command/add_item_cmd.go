package command

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

type AddItemCmd struct {
	UserID            uuid.UUID
	MedicineID        uuid.UUID
	DateOfManufacture time.Time
	Quantity          float32
}

func NewAddItemCmd(userId, medicineId uuid.UUID, date time.Time, quantity float32) (*AddItemCmd, error) {
	if quantity <= 0 {
		return nil, errors.ErrQtyTooLow
	}

	if date.After(time.Now()) {
		return nil, errors.ErrDateTooLate
	}

	return &AddItemCmd{
		UserID:            userId,
		MedicineID:        medicineId,
		DateOfManufacture: date,
		Quantity:          quantity,
	}, nil
}

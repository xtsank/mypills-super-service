package cabinet_item

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

type CabinetItem struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	MedicineID        uuid.UUID
	DateOfManufacture time.Time
	Quantity          float32
}

func NewCabinetItem(id, userID, medID uuid.UUID, date time.Time, quantity float32) (*CabinetItem, error) {
	if quantity <= 0 {
		return nil, errors.ErrQtyTooLow
	}

	if date.After(time.Now()) {
		return nil, errors.ErrDateTooLate
	}

	return &CabinetItem{
		ID:                id,
		UserID:            userID,
		MedicineID:        medID,
		DateOfManufacture: date,
		Quantity:          quantity,
	}, nil
}

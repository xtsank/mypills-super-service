package restock

import (
	"github.com/google/uuid"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
)

type Item struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	MedicineID      uuid.UUID
	MedicineName    string
	DesiredQuantity float64
}

func NewItem(id, userID, medicineID uuid.UUID, medicineName string, desired float64) (*Item, error) {
	if medicineID == uuid.Nil {
		return nil, apperrors.ErrInvalidInput.WithSource()
	}
	item := &Item{ID: id, UserID: userID, MedicineID: medicineID, MedicineName: medicineName}
	if err := item.UpdateDesiredQuantity(desired); err != nil {
		return nil, err
	}
	return item, nil
}

func (i *Item) UpdateDesiredQuantity(desired float64) error {
	if desired <= 0 {
		return apperrors.ErrInvalidInput.WithSource()
	}
	i.DesiredQuantity = desired
	return nil
}

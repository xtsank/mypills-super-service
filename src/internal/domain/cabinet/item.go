package cabinet

import (
	"time"

	"github.com/google/uuid"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
)

type Item struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	CabinetID         uuid.UUID
	MedicineID        uuid.UUID
	MedicineName      string
	ExpireMonths      int
	DateOfManufacture time.Time
	Quantity          float64
	MinimumQuantity   float64
}

func NewItem(id, userID, cabinetID, medicineID uuid.UUID, medicineName string,
	expireMonths int, date time.Time, quantity, minimum float64) (*Item, error) {
	if medicineID == uuid.Nil {
		return nil, apperrors.ErrInvalidInput.WithSource()
	}
	item := &Item{ID: id, UserID: userID, CabinetID: cabinetID, MedicineID: medicineID,
		MedicineName: medicineName, ExpireMonths: expireMonths}
	if err := item.Update(&quantity, &minimum, &date); err != nil {
		return nil, err
	}
	return item, nil
}

func (i *Item) Update(quantity, minimum *float64, date *time.Time) error {
	if quantity != nil && *quantity <= 0 || minimum != nil && *minimum < 0 ||
		date != nil && date.After(time.Now()) {
		return apperrors.ErrInvalidInput.WithSource()
	}
	if quantity != nil {
		i.Quantity = *quantity
	}
	if minimum != nil {
		i.MinimumQuantity = *minimum
	}
	if date != nil {
		i.DateOfManufacture = *date
	}
	return nil
}

func NewUpdatedItem(current *Item, quantity, minimum *float64, date *time.Time) (*Item, error) {
	updated := *current
	if quantity != nil {
		updated.Quantity = *quantity
	}
	if minimum != nil {
		updated.MinimumQuantity = *minimum
	}
	if date != nil {
		updated.DateOfManufacture = *date
	}
	return NewItem(updated.ID, updated.UserID, updated.CabinetID, updated.MedicineID,
		updated.MedicineName, updated.ExpireMonths, updated.DateOfManufacture,
		updated.Quantity, updated.MinimumQuantity)
}

package entity

import (
	"time"

	"github.com/google/uuid"
)

type CabinetItemDetailsEntity struct {
	ID                uuid.UUID `db:"id"`
	UserID            uuid.UUID `db:"user_id"`
	CabinetID         uuid.UUID `db:"cabinet_id"`
	MedicineID        uuid.UUID `db:"medicine_id"`
	MedicineName      string    `db:"medicine_name"`
	ExpireMonths      int       `db:"expire_months"`
	DateOfManufacture time.Time `db:"date_of_manufacture"`
	Quantity          float64   `db:"quantity"`
	MinimumQuantity   float64   `db:"minimum_quantity"`
}

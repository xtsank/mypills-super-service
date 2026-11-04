package entity

import "github.com/google/uuid"

type RestockEntity struct {
	ID              uuid.UUID `db:"id"`
	UserID          uuid.UUID `db:"user_id"`
	MedicineID      uuid.UUID `db:"medicine_id"`
	MedicineName    string    `db:"medicine_name"`
	DesiredQuantity float64   `db:"desired_quantity"`
}

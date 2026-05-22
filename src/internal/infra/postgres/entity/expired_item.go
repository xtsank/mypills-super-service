package entity

import (
	"time"

	"github.com/google/uuid"
)

type ExpiredItemEntity struct {
	ID                uuid.UUID `db:"id"`
	UserID            uuid.UUID `db:"user_id"`
	MedicineID        uuid.UUID `db:"medicine_id"`
	MedicineName      string    `db:"medicine_name"`
	DateOfManufacture time.Time `db:"date_of_manufacture"`
	Quantity          float32   `db:"quantity"`
	ExpiresAt         time.Time `db:"expires_at"`
}


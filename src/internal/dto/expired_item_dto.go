package dto

import (
	"time"

	"github.com/google/uuid"
)

type ExpiredItemDto struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	MedicineID        uuid.UUID
	MedicineName      string
	DateOfManufacture time.Time
	Quantity          float32
	ExpiresAt         time.Time
}


package entity

import (
	"time"

	"github.com/google/uuid"
)

type CabinetItemEntity struct {
	ID                uuid.UUID `bson:"_id"`
	UserID            uuid.UUID `bson:"user_id"`
	MedicineID        uuid.UUID `bson:"medicine_id"`
	DateOfManufacture time.Time `bson:"date_of_manufacture"`
	Quantity          float32   `bson:"quantity"`
}

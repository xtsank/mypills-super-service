package req

import (
	"time"

	"github.com/google/uuid"
)

type AddItemDto struct {
	MedicineID        uuid.UUID `json:"medicine_id" binding:"required"`
	DateOfManufacture time.Time `json:"date_of_manufacture" binding:"required"`
	Quantity          float32   `json:"quantity" binding:"required"`
}

package req

import "github.com/google/uuid"

type CabinetItemCreateDto struct {
	MedicineID        uuid.UUID `json:"medicine_id" required:"true"`
	DateOfManufacture DateDto   `json:"date_of_manufacture" required:"true"`
	Quantity          float64   `json:"quantity" required:"true"`
	MinimumQuantity   float64   `json:"minimum_quantity" required:"true"`
}

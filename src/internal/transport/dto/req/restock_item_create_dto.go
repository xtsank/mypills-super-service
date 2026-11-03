package req

import "github.com/google/uuid"

type RestockItemCreateDto struct {
	MedicineID      uuid.UUID `json:"medicine_id" required:"true"`
	DesiredQuantity float64   `json:"desired_quantity" required:"true"`
}

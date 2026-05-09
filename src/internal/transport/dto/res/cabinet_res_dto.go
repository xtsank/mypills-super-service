package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
)

type CabinetResDto struct {
	ID         uuid.UUID `json:"id"`
	MedicineID uuid.UUID `json:"medicine_id"`
	Quantity   float32   `json:"quantity"`
}

func NewCabinetResDto(item *cabinet_item.CabinetItem) *CabinetResDto {
	return &CabinetResDto{
		ID:         item.ID,
		MedicineID: item.MedicineID,
		Quantity:   item.Quantity,
	}
}

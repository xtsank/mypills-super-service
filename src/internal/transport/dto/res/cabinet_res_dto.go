package res

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

type CabinetResDto struct {
	ID         uuid.UUID `json:"id"`
	MedicineID uuid.UUID `json:"medicine_id"`
	Quantity   float32   `json:"quantity"`
}

type CabinetItemDetailsResDto struct {
	ID                uuid.UUID `json:"id"`
	MedicineID        uuid.UUID `json:"medicine_id"`
	MedicineName      string    `json:"medicine_name"`
	DateOfManufacture time.Time `json:"date_of_manufacture"`
	ExpiresAt         time.Time `json:"expires_at"`
	Quantity          float32   `json:"quantity"`
}

func NewCabinetResDto(item *cabinet_item.CabinetItem) *CabinetResDto {
	return &CabinetResDto{
		ID:         item.ID,
		MedicineID: item.MedicineID,
		Quantity:   item.Quantity,
	}
}

func NewCabinetItemDetailsResDto(items []entity.CabinetItemWithMedicineEntity) []CabinetItemDetailsResDto {
	res := make([]CabinetItemDetailsResDto, 0, len(items))
	for _, item := range items {
		res = append(res, CabinetItemDetailsResDto{
			ID:                item.ID,
			MedicineID:        item.MedicineID,
			MedicineName:      item.MedicineName,
			DateOfManufacture: item.DateOfManufacture,
			ExpiresAt:         item.ExpiresAt,
			Quantity:          item.Quantity,
		})
	}
	return res
}

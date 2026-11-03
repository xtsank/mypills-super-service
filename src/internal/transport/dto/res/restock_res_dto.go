package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/restock"
)

type RestockItemDto struct {
	ID              uuid.UUID `json:"id"`
	MedicineID      uuid.UUID `json:"medicine_id"`
	MedicineName    string    `json:"medicine_name"`
	DesiredQuantity float64   `json:"desired_quantity"`
}

type RestockPageDto struct {
	Items  []RestockItemDto `json:"items"`
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}

func NewRestockItemDto(item restock.Item) *RestockItemDto {
	return &RestockItemDto{
		ID:              item.ID,
		MedicineID:      item.MedicineID,
		MedicineName:    item.MedicineName,
		DesiredQuantity: item.DesiredQuantity,
	}
}

func NewRestockPageDto(page restock.Page) *RestockPageDto {
	result := &RestockPageDto{
		Items:  make([]RestockItemDto, 0, len(page.Items)),
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
	for _, item := range page.Items {
		result.Items = append(result.Items, *NewRestockItemDto(item))
	}
	return result
}

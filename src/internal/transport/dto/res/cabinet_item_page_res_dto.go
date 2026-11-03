package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
)

type CabinetItemDto struct {
	ID                uuid.UUID `json:"id"`
	CabinetID         uuid.UUID `json:"cabinet_id"`
	MedicineID        uuid.UUID `json:"medicine_id"`
	MedicineName      string    `json:"medicine_name"`
	DateOfManufacture string    `json:"date_of_manufacture"`
	ExpiresAt         string    `json:"expires_at"`
	Quantity          float64   `json:"quantity"`
	MinimumQuantity   float64   `json:"minimum_quantity"`
	LowStock          bool      `json:"low_stock"`
}

type CabinetItemsPageDto struct {
	Items  []CabinetItemDto `json:"items"`
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}

func NewCabinetItemDto(item cabinet.Item) *CabinetItemDto {
	return &CabinetItemDto{
		ID:                item.ID,
		CabinetID:         item.CabinetID,
		MedicineID:        item.MedicineID,
		MedicineName:      item.MedicineName,
		DateOfManufacture: item.DateOfManufacture.Format("2006-01-02"),
		ExpiresAt:         item.DateOfManufacture.AddDate(0, item.ExpireMonths, 0).Format("2006-01-02"),
		Quantity:          item.Quantity,
		MinimumQuantity:   item.MinimumQuantity,
		LowStock:          item.Quantity <= item.MinimumQuantity,
	}
}

func NewCabinetItemsPageDto(page cabinet.ItemPage) *CabinetItemsPageDto {
	result := &CabinetItemsPageDto{
		Items: make([]CabinetItemDto, 0, len(page.Items)),
		Total: page.Total, Limit: page.Limit, Offset: page.Offset,
	}
	for _, item := range page.Items {
		result.Items = append(result.Items, *NewCabinetItemDto(item))
	}
	return result
}

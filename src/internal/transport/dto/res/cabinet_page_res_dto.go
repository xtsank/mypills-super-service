package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
)

type CabinetDto struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CabinetsPageDto struct {
	Items  []CabinetDto `json:"items"`
	Total  int          `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

func NewCabinetDto(item cabinet.Cabinet) *CabinetDto {
	return &CabinetDto{ID: item.ID, Name: item.Name}
}

func NewCabinetsPageDto(page cabinet.Page) *CabinetsPageDto {
	result := &CabinetsPageDto{
		Items:  make([]CabinetDto, 0, len(page.Items)),
		Total:  page.Total,
		Limit:  page.Limit,
		Offset: page.Offset,
	}
	for _, item := range page.Items {
		result.Items = append(result.Items, *NewCabinetDto(item))
	}
	return result
}

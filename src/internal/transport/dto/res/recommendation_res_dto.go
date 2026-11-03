package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
)

type RecommendationDto struct {
	MedicineID          uuid.UUID `json:"medicine_id"`
	CabinetID           uuid.UUID `json:"cabinet_id"`
	CabinetItemID       uuid.UUID `json:"cabinet_item_id"`
	Name                string    `json:"name"`
	MethodOfApplication string    `json:"method_of_application"`
	Dosage              float32   `json:"dosage"`
	UnitName            string    `json:"unit_name"`
	FrequencyPerDay     int       `json:"frequency_per_day"`
	QuantityInCabinet   float64   `json:"quantity_in_cabinet"`
}

type RecommendationResultDto struct {
	Items []RecommendationDto `json:"items"`
}

func NewRecommendationResultDto(items []medicine.Recommendation) *RecommendationResultDto {
	result := &RecommendationResultDto{Items: make([]RecommendationDto, 0, len(items))}
	for _, item := range items {
		result.Items = append(result.Items, RecommendationDto{
			MedicineID: item.MedicineID, CabinetID: item.CabinetID, CabinetItemID: item.CabinetItemID,
			Name: item.Name, MethodOfApplication: item.MethodOfApplication,
			Dosage: item.Dosage, UnitName: item.UnitName, FrequencyPerDay: item.FrequencyPerDay,
			QuantityInCabinet: item.QuantityInCabinet,
		})
	}
	return result
}

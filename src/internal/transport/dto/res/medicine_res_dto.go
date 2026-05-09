package res

import "github.com/google/uuid"

type MedicineRecommendation struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	MethodOfApplication string    `json:"method_of_application"`
	Dosage              float32   `json:"dosage"`
	UnitName            string    `json:"unit_name"`
	Frequency           int       `json:"frequency"`
	QuantityInCabinet   float32   `json:"quantity_in_cabinet"`
}

type MedicineResDto struct {
	Recommendations []*MedicineRecommendation `json:"recommendations"`
}

func NewMedicineResDto(recs []*MedicineRecommendation) *MedicineResDto {
	return &MedicineResDto{
		Recommendations: recs,
	}
}

package medicine

import "github.com/google/uuid"

type Recommendation struct {
	MedicineID          uuid.UUID
	CabinetID           uuid.UUID
	CabinetItemID       uuid.UUID
	Name                string
	MethodOfApplication string
	Dosage              float32
	UnitName            string
	FrequencyPerDay     int
	QuantityInCabinet   float64
}

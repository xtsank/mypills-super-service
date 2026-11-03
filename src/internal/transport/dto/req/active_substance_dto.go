package req

import "github.com/google/uuid"

type ActiveSubstanceDto struct {
	ID            uuid.UUID `json:"id"`
	Concentration float32   `json:"concentration"`
}

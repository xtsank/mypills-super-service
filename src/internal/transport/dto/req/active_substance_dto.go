package req

import "github.com/google/uuid"

type ActiveSubstanceDto struct {
	ID            uuid.UUID
	Concentration float32
}

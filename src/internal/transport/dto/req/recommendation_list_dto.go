package req

import "github.com/google/uuid"

type RecommendationListDto struct {
	IllnessID uuid.UUID
	CabinetID *uuid.UUID
}

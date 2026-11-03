package command

import "github.com/google/uuid"

type ListRecommendationsCmd struct {
	UserID    uuid.UUID
	IllnessID uuid.UUID
	CabinetID *uuid.UUID
}

func NewListRecommendationsCmd(userID, illnessID uuid.UUID, cabinetID *uuid.UUID) *ListRecommendationsCmd {
	return &ListRecommendationsCmd{UserID: userID, IllnessID: illnessID, CabinetID: cabinetID}
}

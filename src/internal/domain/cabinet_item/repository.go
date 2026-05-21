package cabinet_item

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/dto"
)

type ICabinetItemRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*CabinetItem, error)
	FindExistingCabinetItem(ctx context.Context, userID uuid.UUID, medicineID uuid.UUID, date time.Time) (*CabinetItem, error)
	FindById(ctx context.Context, id uuid.UUID) (*CabinetItem, error)
	FindExpiredByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.ExpiredItemDto, error)
	Update(ctx context.Context, cabinetItem *CabinetItem) error
	Save(ctx context.Context, cabinetItem *CabinetItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

package cabinet_item

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/dto"
<<<<<<< HEAD
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
=======
>>>>>>> 1f83dea7bd71d6b52bdd54933e14f6e23c6bc04a
)

type ICabinetItemRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*CabinetItem, error)
	FindDetailsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.CabinetItemWithMedicineEntity, error)
	FindExistingCabinetItem(ctx context.Context, userID uuid.UUID, medicineID uuid.UUID, date time.Time) (*CabinetItem, error)
	FindById(ctx context.Context, id uuid.UUID) (*CabinetItem, error)
	FindExpiredByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.ExpiredItemDto, error)
	Update(ctx context.Context, cabinetItem *CabinetItem) error
	Save(ctx context.Context, cabinetItem *CabinetItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

package medicine

import (
	"context"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

type IMedicineRepository interface {
	FindByIllness(ctx context.Context, illnessID uuid.UUID) ([]*Medicine, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Medicine, error)
	ListDosageRules(ctx context.Context) ([]entity.DosageRuleWithMedicineEntity, error)

	Create(ctx context.Context, medicine *Medicine) error
	Update(ctx context.Context, medicine *Medicine) error
	Delete(ctx context.Context, id uuid.UUID) error

	UpdateIndications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error
	UpdateContraindications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error
	UpdateComposition(ctx context.Context, medicineID uuid.UUID, substances []ActiveSubstance) error
	AddDosageRule(ctx context.Context, medicineID uuid.UUID, rule *DosageRule) error
	DeleteDosageRule(ctx context.Context, medicineID uuid.UUID) error
}

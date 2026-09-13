package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

type MockMedicineRepository struct {
	mock.Mock
}

func (m *MockMedicineRepository) FindByIllness(ctx context.Context, illnessID uuid.UUID) ([]*medicine.Medicine, error) {
	args := m.Called(ctx, illnessID)
	result, _ := args.Get(0).([]*medicine.Medicine)
	return result, args.Error(1)
}

func (m *MockMedicineRepository) FindByID(ctx context.Context, id uuid.UUID) (*medicine.Medicine, error) {
	args := m.Called(ctx, id)
	result, _ := args.Get(0).(*medicine.Medicine)
	return result, args.Error(1)
}

func (m *MockMedicineRepository) ListDosageRules(ctx context.Context) ([]entity.DosageRuleWithMedicineEntity, error) {
	args := m.Called(ctx)
	result, _ := args.Get(0).([]entity.DosageRuleWithMedicineEntity)
	return result, args.Error(1)
}

func (m *MockMedicineRepository) Create(ctx context.Context, value *medicine.Medicine) error {
	return m.Called(ctx, value).Error(0)
}

func (m *MockMedicineRepository) Update(ctx context.Context, value *medicine.Medicine) error {
	return m.Called(ctx, value).Error(0)
}

func (m *MockMedicineRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockMedicineRepository) UpdateIndications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error {
	return m.Called(ctx, medicineID, ids).Error(0)
}

func (m *MockMedicineRepository) UpdateContraindications(ctx context.Context, medicineID uuid.UUID, ids []uuid.UUID) error {
	return m.Called(ctx, medicineID, ids).Error(0)
}

func (m *MockMedicineRepository) UpdateComposition(ctx context.Context, medicineID uuid.UUID, substances []medicine.ActiveSubstance) error {
	return m.Called(ctx, medicineID, substances).Error(0)
}

func (m *MockMedicineRepository) AddDosageRule(ctx context.Context, medicineID uuid.UUID, rule *medicine.DosageRule) error {
	return m.Called(ctx, medicineID, rule).Error(0)
}

func (m *MockMedicineRepository) DeleteDosageRule(ctx context.Context, ruleID uuid.UUID) error {
	return m.Called(ctx, ruleID).Error(0)
}

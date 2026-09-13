package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/dto"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

type MockCabinetRepository struct {
	mock.Mock
}

func (m *MockCabinetRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*cabinet_item.CabinetItem, error) {
	args := m.Called(ctx, userID)

	var result []*cabinet_item.CabinetItem
	if args.Get(0) != nil {
		result = args.Get(0).([]*cabinet_item.CabinetItem)
	}

	return result, args.Error(1)
}

func (m *MockCabinetRepository) FindDetailsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.CabinetItemWithMedicineEntity, error) {
	args := m.Called(ctx, userID)

	var result []entity.CabinetItemWithMedicineEntity
	if args.Get(0) != nil {
		result = args.Get(0).([]entity.CabinetItemWithMedicineEntity)
	}

	return result, args.Error(1)
}

func (m *MockCabinetRepository) FindExistingCabinetItem(ctx context.Context, userID uuid.UUID, medicineID uuid.UUID, date time.Time) (*cabinet_item.CabinetItem, error) {
	args := m.Called(ctx, userID, medicineID, date)

	var result *cabinet_item.CabinetItem
	if args.Get(0) != nil {
		result = args.Get(0).(*cabinet_item.CabinetItem)
	}

	return result, args.Error(1)
}

func (m *MockCabinetRepository) FindById(ctx context.Context, id uuid.UUID) (*cabinet_item.CabinetItem, error) {
	args := m.Called(ctx, id)

	var result *cabinet_item.CabinetItem
	if args.Get(0) != nil {
		result = args.Get(0).(*cabinet_item.CabinetItem)
	}

	return result, args.Error(1)
}

func (m *MockCabinetRepository) FindExpiredByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.ExpiredItemDto, error) {
	args := m.Called(ctx, userID)

	var result []*dto.ExpiredItemDto
	if args.Get(0) != nil {
		result = args.Get(0).([]*dto.ExpiredItemDto)
	}

	return result, args.Error(1)
}

func (m *MockCabinetRepository) Update(ctx context.Context, cabinetItem *cabinet_item.CabinetItem) error {
	return m.Called(ctx, cabinetItem).Error(0)
}

func (m *MockCabinetRepository) Save(ctx context.Context, cabinetItem *cabinet_item.CabinetItem) error {
	return m.Called(ctx, cabinetItem).Error(0)
}

func (m *MockCabinetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

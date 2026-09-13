package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	args := m.Called(ctx, login)
	result, _ := args.Get(0).(*user.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	result, _ := args.Get(0).(*user.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) FindNotifyEnabled(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	result, _ := args.Get(0).([]*user.User)
	return result, args.Error(1)
}

func (m *MockUserRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	args := m.Called(ctx, login)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, value *user.User) error {
	return m.Called(ctx, value).Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, value *user.User) error {
	return m.Called(ctx, value).Error(0)
}

func (m *MockUserRepository) UpdateNotify(ctx context.Context, id uuid.UUID, enabled bool, intervalMinutes int, lastNotifiedAt *time.Time) error {
	return m.Called(ctx, id, enabled, intervalMinutes, lastNotifiedAt).Error(0)
}

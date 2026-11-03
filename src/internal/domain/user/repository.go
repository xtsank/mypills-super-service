package user

import (
	"context"

	"github.com/google/uuid"
)

type IUserRepository interface {
	FindByLogin(ctx context.Context, login string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	ExistsByLogin(ctx context.Context, login string) (bool, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

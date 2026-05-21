package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type IUserRepository interface {
	FindByLogin(ctx context.Context, login string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindNotifyEnabled(ctx context.Context) ([]*User, error)
	ExistsByLogin(ctx context.Context, login string) (bool, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateNotify(ctx context.Context, id uuid.UUID, enabled bool, intervalMinutes int, lastNotifiedAt *time.Time) error
}

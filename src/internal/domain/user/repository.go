package user

import (
	"context"
)

type IUserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByLogin(ctx context.Context, login string) (*User, error)
}

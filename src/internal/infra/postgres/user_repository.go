package postgres

import (
	"context"

	"github.com/xtsank/mypills-super-service/internal/domain/user"
)

type PostgresUserRepository struct {
	db string
}

func NewPostgresUserRepository(db string) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *user.User) error {
	return nil
}

func (r *PostgresUserRepository) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	return nil, nil
}

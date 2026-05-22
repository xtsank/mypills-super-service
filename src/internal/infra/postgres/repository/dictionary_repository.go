package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/dictionary"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
)

type PostgresDictionaryRepository struct {
	db *sqlx.DB
}

func NewPostgresDictionaryRepository(i do.Injector) (dictionary.IRepository, error) {
	db := do.MustInvoke[*sqlx.DB](i)
	return &PostgresDictionaryRepository{db: db}, nil
}

func (r *PostgresDictionaryRepository) ListIllnesses(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, queries.Dictionary.ListIllnesses)
}

func (r *PostgresDictionaryRepository) ListSubstances(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, queries.Dictionary.ListSubstances)
}

func (r *PostgresDictionaryRepository) ListForms(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, queries.Dictionary.ListForms)
}

func (r *PostgresDictionaryRepository) ListUnits(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, queries.Dictionary.ListUnits)
}

func (r *PostgresDictionaryRepository) ListMedicines(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, queries.Dictionary.ListMedicines)
}

func (r *PostgresDictionaryRepository) list(ctx context.Context, query string) ([]dictionary.Item, error) {
	var rows []struct {
		ID   uuid.UUID `db:"id"`
		Name string    `db:"name"`
	}

	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	items := make([]dictionary.Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, dictionary.Item{ID: row.ID, Name: row.Name})
	}
	return items, nil
}



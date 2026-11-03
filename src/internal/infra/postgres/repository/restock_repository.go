package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/restock"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
)

type PostgresRestockRepository struct {
	db *sqlx.DB
}

func NewPostgresRestockRepository(i do.Injector) (restock.IRestockRepository, error) {
	return &PostgresRestockRepository{db: do.MustInvoke[*sqlx.DB](i)}, nil
}

func mapRestockEntity(e entity.RestockEntity) restock.Item {
	return restock.Item{ID: e.ID, UserID: e.UserID, MedicineID: e.MedicineID,
		MedicineName: e.MedicineName, DesiredQuantity: e.DesiredQuantity}
}

func (r *PostgresRestockRepository) List(ctx context.Context, userID uuid.UUID, id *uuid.UUID, limit, offset int) (restock.Page, error) {
	p := restock.Page{Items: []restock.Item{}, Limit: limit, Offset: offset}
	var idFilter any
	if id != nil {
		idFilter = *id
	}
	if err := r.db.GetContext(ctx, &p.Total, queries.Restock.Count, userID, idFilter); err != nil {
		return p, apperrors.ErrInternal.WithError(err)
	}
	var entities []entity.RestockEntity
	if err := r.db.SelectContext(ctx, &entities, queries.Restock.List, userID, idFilter, limit, offset); err != nil {
		return p, apperrors.ErrInternal.WithError(err)
	}
	for _, e := range entities {
		p.Items = append(p.Items, mapRestockEntity(e))
	}
	return p, nil
}

func (r *PostgresRestockRepository) Find(ctx context.Context, userID, id uuid.UUID) (*restock.Item, error) {
	var e entity.RestockEntity
	err := r.db.GetContext(ctx, &e, queries.Restock.Find, userID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.ErrInternal.WithError(err)
	}
	item := mapRestockEntity(e)
	return &item, nil
}

func (r *PostgresRestockRepository) Create(ctx context.Context, item restock.Item) error {
	_, err := r.db.ExecContext(ctx, queries.Restock.Create,
		item.ID, item.UserID, item.MedicineID, item.DesiredQuantity)
	if err != nil {
		var pgErr interface{ SQLState() string }
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return apperrors.ErrRestockItemConflict.WithError(err)
		}
		return apperrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresRestockRepository) Update(ctx context.Context, item restock.Item) error {
	result, err := r.db.ExecContext(ctx, queries.Restock.Update,
		item.ID, item.UserID, item.DesiredQuantity)
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	if count == 0 {
		return apperrors.ErrRestockItemNotFound.WithSource()
	}
	return nil
}

func (r *PostgresRestockRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, queries.Restock.Delete, id, userID)
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	if count == 0 {
		return apperrors.ErrRestockItemNotFound.WithSource()
	}
	return nil
}

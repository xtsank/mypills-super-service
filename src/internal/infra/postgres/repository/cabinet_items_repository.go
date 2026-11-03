package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
)

type PostgresCabinetItemRepository struct {
	db *sqlx.DB
}

func NewPostgresCabinetItemRepository(i do.Injector) (cabinet.ICabinetItemRepository, error) {
	return &PostgresCabinetItemRepository{db: do.MustInvoke[*sqlx.DB](i)}, nil
}

func mapCabinetItemEntity(e entity.CabinetItemDetailsEntity) cabinet.Item {
	return cabinet.Item{
		ID:                e.ID,
		UserID:            e.UserID,
		CabinetID:         e.CabinetID,
		MedicineID:        e.MedicineID,
		MedicineName:      e.MedicineName,
		ExpireMonths:      e.ExpireMonths,
		DateOfManufacture: e.DateOfManufacture,
		Quantity:          e.Quantity,
		MinimumQuantity:   e.MinimumQuantity,
	}
}

func (r *PostgresCabinetItemRepository) List(ctx context.Context, userID, cabinetID uuid.UUID) ([]cabinet.Item, error) {
	var entities []entity.CabinetItemDetailsEntity
	if err := r.db.SelectContext(ctx, &entities, queries.CabinetItems.List, userID, cabinetID); err != nil {
		return nil, apperrors.ErrInternal.WithError(err)
	}
	items := make([]cabinet.Item, 0, len(entities))
	for _, entity := range entities {
		items = append(items, mapCabinetItemEntity(entity))
	}
	return items, nil
}

func (r *PostgresCabinetItemRepository) Find(ctx context.Context, userID, cabinetID, id uuid.UUID) (*cabinet.Item, error) {
	var result entity.CabinetItemDetailsEntity
	err := r.db.GetContext(ctx, &result, queries.CabinetItems.Find, userID, cabinetID, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, apperrors.ErrInternal.WithError(err)
	}
	item := mapCabinetItemEntity(result)
	return &item, nil
}

func (r *PostgresCabinetItemRepository) Create(ctx context.Context, item cabinet.Item) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	defer tx.Rollback()
	key := fmt.Sprintf("%s:%s:%s", item.CabinetID, item.MedicineID, item.DateOfManufacture.Format("2006-01-02"))
	if _, err := tx.ExecContext(ctx, queries.CabinetItems.LockDuplicate, key); err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	var exists bool
	if err := tx.GetContext(ctx, &exists, queries.CabinetItems.ExistsDuplicate,
		item.CabinetID, item.MedicineID, item.DateOfManufacture); err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	if exists {
		return apperrors.ErrCabinetItemConflict.WithSource()
	}
	_, err = tx.ExecContext(ctx, queries.CabinetItems.Create, item.ID, item.UserID, item.CabinetID, item.MedicineID,
		item.DateOfManufacture, item.Quantity, item.MinimumQuantity)
	if err != nil {
		return mapCabinetDBError(err)
	}
	if err := tx.Commit(); err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresCabinetItemRepository) Update(ctx context.Context, item cabinet.Item) error {
	result, err := r.db.ExecContext(ctx, queries.CabinetItems.Update,
		item.ID, item.UserID, item.CabinetID, item.Quantity, item.MinimumQuantity,
		item.DateOfManufacture, item.MedicineID)
	if err != nil {
		return mapCabinetDBError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	if count == 0 {
		return apperrors.ErrCabinetItemConflict.WithSource()
	}
	return nil
}

func (r *PostgresCabinetItemRepository) Delete(ctx context.Context, userID, cabinetID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, queries.CabinetItems.Delete, id, userID, cabinetID)
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	if count == 0 {
		return apperrors.ErrCabinetItemNotFound.WithSource()
	}
	return nil
}

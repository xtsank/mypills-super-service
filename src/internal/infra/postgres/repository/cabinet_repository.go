package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
)

type PostgresCabinetRepository struct {
	db *sqlx.DB
}

func NewPostgresCabinetRepository(i do.Injector) (cabinet.ICabinetRepository, error) {
	return &PostgresCabinetRepository{db: do.MustInvoke[*sqlx.DB](i)}, nil
}

func (r *PostgresCabinetRepository) EnsureLegacy(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, queries.Cabinet.EnsureLegacy, userID)
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	_, err = r.db.ExecContext(ctx, queries.Cabinet.AssignLegacyItems, userID)
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresCabinetRepository) List(ctx context.Context, userID uuid.UUID, id *uuid.UUID, limit, offset int) (cabinet.Page, error) {
	p := cabinet.Page{Items: []cabinet.Cabinet{}, Limit: limit, Offset: offset}
	var idFilter any
	if id != nil {
		idFilter = *id
	}
	if err := r.db.GetContext(ctx, &p.Total, queries.Cabinet.Count, userID, idFilter); err != nil {
		return p, apperrors.ErrInternal.WithError(err)
	}
	var entities []entity.CabinetEntity
	if err := r.db.SelectContext(ctx, &entities, queries.Cabinet.List, userID, idFilter, limit, offset); err != nil {
		return p, apperrors.ErrInternal.WithError(err)
	}
	for _, item := range entities {
		p.Items = append(p.Items, cabinet.Cabinet{ID: item.ID, UserID: item.UserID, Name: item.Name})
	}
	return p, nil
}

func (r *PostgresCabinetRepository) Find(ctx context.Context, userID, id uuid.UUID) (*cabinet.Cabinet, error) {
	var item entity.CabinetEntity
	if err := r.db.GetContext(ctx, &item, queries.Cabinet.Find, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, apperrors.ErrInternal.WithError(err)
	}
	return &cabinet.Cabinet{ID: item.ID, UserID: item.UserID, Name: item.Name}, nil
}

func (r *PostgresCabinetRepository) Create(ctx context.Context, item cabinet.Cabinet) error {
	_, err := r.db.ExecContext(ctx, queries.Cabinet.Create, item.ID, item.UserID, item.Name)
	return mapCabinetDBError(err)
}

func (r *PostgresCabinetRepository) Rename(ctx context.Context, item cabinet.Cabinet) error {
	result, err := r.db.ExecContext(ctx, queries.Cabinet.Rename, item.ID, item.UserID, item.Name)
	if err != nil {
		return mapCabinetDBError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	if count == 0 {
		return apperrors.ErrCabinetNotFound.WithSource()
	}
	return nil
}

func (r *PostgresCabinetRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, queries.Cabinet.Delete, id, userID)
	if err != nil {
		return mapCabinetDBError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return apperrors.ErrInternal.WithError(err)
	}
	if count == 0 {
		return apperrors.ErrCabinetConflict.WithSource()
	}
	return nil
}

func mapCabinetDBError(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && (pgErr.Code == "23505" || pgErr.Code == "23503") {
		return apperrors.ErrCabinetConflict.WithError(err)
	}
	return apperrors.ErrInternal.WithError(err)
}

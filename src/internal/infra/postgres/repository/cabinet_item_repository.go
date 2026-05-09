package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

type PostgresCabinetItemRepository struct {
	db *sqlx.DB
}

func NewPostgresCabinetItemRepository(i do.Injector) (cabinet_item.ICabinetItemRepository, error) {
	db := do.MustInvoke[*sqlx.DB](i)
	return &PostgresCabinetItemRepository{db: db}, nil
}

func (r *PostgresCabinetItemRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*cabinet_item.CabinetItem, error) {
	var ents []entity.CabinetItemEntity
	query := `SELECT * FROM User_Medicine WHERE user_id = $1`

	err := r.db.SelectContext(ctx, &ents, query, userID)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	res := make([]*cabinet_item.CabinetItem, len(ents))
	for i, ent := range ents {
		res[i], err = cabinet_item.NewCabinetItem(ent.ID, ent.UserID, ent.MedicineID, ent.DateOfManufacture, ent.Quantity)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (r *PostgresCabinetItemRepository) FindExistingCabinetItem(ctx context.Context, userID uuid.UUID, medID uuid.UUID, date time.Time) (*cabinet_item.CabinetItem, error) {
	var ent entity.CabinetItemEntity
	query := `SELECT * FROM User_Medicine 
              WHERE user_id = $1 AND medicine_id = $2 AND date_of_manufacture = $3::date`

	err := r.db.GetContext(ctx, &ent, query, userID, medID, date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return cabinet_item.NewCabinetItem(ent.ID, ent.UserID, ent.MedicineID, ent.DateOfManufacture, ent.Quantity)
}

func (r *PostgresCabinetItemRepository) FindById(ctx context.Context, id uuid.UUID) (*cabinet_item.CabinetItem, error) {
	var ent entity.CabinetItemEntity
	query := `SELECT * FROM User_Medicine WHERE id = $1`

	err := r.db.GetContext(ctx, &ent, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return cabinet_item.NewCabinetItem(ent.ID, ent.UserID, ent.MedicineID, ent.DateOfManufacture, ent.Quantity)
}

func (r *PostgresCabinetItemRepository) Update(ctx context.Context, item *cabinet_item.CabinetItem) error {
	ent := entity.CabinetItemEntity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}

	query := `UPDATE User_Medicine SET quantity = :quantity WHERE id = :id`

	_, err := r.db.NamedExecContext(ctx, query, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresCabinetItemRepository) Save(ctx context.Context, item *cabinet_item.CabinetItem) error {
	ent := entity.CabinetItemEntity{
		ID:                item.ID,
		UserID:            item.UserID,
		MedicineID:        item.MedicineID,
		DateOfManufacture: item.DateOfManufacture,
		Quantity:          item.Quantity,
	}

	query := `INSERT INTO User_Medicine (id, user_id, medicine_id, date_of_manufacture, quantity)
              VALUES (:id, :user_id, :medicine_id, :date_of_manufacture, :quantity)`

	_, err := r.db.NamedExecContext(ctx, query, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *PostgresCabinetItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM User_Medicine WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

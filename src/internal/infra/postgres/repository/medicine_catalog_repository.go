package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
)

func NewPostgresMedicineCatalogRepository(i do.Injector) (medicine.ICatalogRepository, error) {
	return &PostgresMedicineRepository{db: do.MustInvoke[*sqlx.DB](i)}, nil
}

func (r *PostgresMedicineRepository) List(ctx context.Context, id *uuid.UUID, name string, limit, offset int) (medicine.Page, error) {
	page := medicine.Page{Items: []*medicine.Medicine{}, Limit: limit, Offset: offset}
	name = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(name)
	var idFilter any
	if id != nil {
		idFilter = *id
	}
	if err := r.db.GetContext(ctx, &page.Total, queries.Medicine.Count, idFilter, name); err != nil {
		return page, apperrors.ErrInternal.WithError(err)
	}
	var ids []uuid.UUID
	if err := r.db.SelectContext(ctx, &ids, queries.Medicine.ListIDs, idFilter, name, limit, offset); err != nil {
		return page, apperrors.ErrInternal.WithError(err)
	}
	for _, medicineID := range ids {
		item, err := r.FindByID(ctx, medicineID)
		if err != nil {
			return page, err
		}
		page.Items = append(page.Items, item)
	}
	return page, nil
}

package lab3

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

type SQLCatalog struct {
	db *sqlx.DB
}

func NewSQLCatalog(pool *sql.DB) *SQLCatalog {
	return &SQLCatalog{db: sqlx.NewDb(pool, "pgx")}
}

func (c *SQLCatalog) Find(ctx context.Context, id uuid.UUID) (entity.MedicineEntity, error) {
	var medicine entity.MedicineEntity
	err := c.db.GetContext(ctx, &medicine,
		"SELECT "+medicineFields+" FROM medicine WHERE id = $1 LIMIT $2", id, 1)
	return medicine, err
}

func (c *SQLCatalog) Count(ctx context.Context, name string) (int64, error) {
	var count int64
	err := c.db.GetContext(ctx, &count,
		"SELECT count(*) FROM medicine WHERE name ILIKE $1 ESCAPE chr(92)", namePattern(name))
	return count, err
}

func (c *SQLCatalog) Page(ctx context.Context, name string, offset int) ([]entity.MedicineEntity, error) {
	items := []entity.MedicineEntity{}
	query := "SELECT " + medicineFields + " FROM medicine WHERE name ILIKE $1 ESCAPE chr(92)" +
		" ORDER BY name, id LIMIT $2"
	args := []any{namePattern(name), pageSize}
	if offset > 0 {
		query += " OFFSET $3"
		args = append(args, offset)
	}
	err := c.db.SelectContext(ctx, &items, query, args...)
	return items, err
}

func (c *SQLCatalog) Insert(ctx context.Context, medicine entity.MedicineEntity) error {
	return c.write(ctx, "INSERT INTO medicine ("+medicineFields+") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		medicine.ID, medicine.FormId, medicine.UnitId, medicine.Name, medicine.ExpireTime,
		medicine.EffectOnDriver, medicine.EffectOnPregnant, medicine.MethodOfApplication, medicine.IsPrescription)
}

func (c *SQLCatalog) Rename(ctx context.Context, id uuid.UUID, name string) error {
	return c.write(ctx, "UPDATE medicine SET name = $1 WHERE id = $2", name, id)
}

func (c *SQLCatalog) Delete(ctx context.Context, id uuid.UUID) error {
	return c.write(ctx, "DELETE FROM medicine WHERE id = $1", id)
}

func (c *SQLCatalog) write(ctx context.Context, query string, args ...any) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if err := requireOneRow(rows); err != nil {
		return err
	}
	return tx.Commit()
}

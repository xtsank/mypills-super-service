package lab3

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GORMCatalog struct {
	db *gorm.DB
}

func NewGORMCatalog(pool *sql.DB) (*GORMCatalog, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: pool}), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	return &GORMCatalog{db: db}, nil
}

func (c *GORMCatalog) Find(ctx context.Context, id uuid.UUID) (entity.MedicineEntity, error) {
	var medicine entity.MedicineEntity
	err := c.db.WithContext(ctx).Table("medicine").Select(medicineFields).
		Where("id = ?", id).Take(&medicine).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = sql.ErrNoRows
	}
	return medicine, err
}

func (c *GORMCatalog) Count(ctx context.Context, name string) (int64, error) {
	var count int64
	err := c.db.WithContext(ctx).Table("medicine").
		Where("name ILIKE ? ESCAPE chr(92)", namePattern(name)).Count(&count).Error
	return count, err
}

func (c *GORMCatalog) Page(ctx context.Context, name string, offset int) ([]entity.MedicineEntity, error) {
	items := []entity.MedicineEntity{}
	err := c.db.WithContext(ctx).Table("medicine").Select(medicineFields).
		Where("name ILIKE ? ESCAPE chr(92)", namePattern(name)).
		Order("name, id").Limit(pageSize).Offset(offset).Find(&items).Error
	return items, err
}

func (c *GORMCatalog) Insert(ctx context.Context, medicine entity.MedicineEntity) error {
	return c.write(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Create(&medicine)
	})
}

func (c *GORMCatalog) Rename(ctx context.Context, id uuid.UUID, name string) error {
	return c.write(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&entity.MedicineEntity{}).Where("id = ?", id).Update("name", name)
	})
}

func (c *GORMCatalog) Delete(ctx context.Context, id uuid.UUID) error {
	return c.write(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.Where("id = ?", id).Delete(&entity.MedicineEntity{})
	})
}

func (c *GORMCatalog) write(ctx context.Context, operation func(*gorm.DB) *gorm.DB) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := operation(tx.Table("medicine"))
		if result.Error != nil {
			return result.Error
		}
		return requireOneRow(result.RowsAffected)
	})
}

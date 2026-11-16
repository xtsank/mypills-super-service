package lab3

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

func verifyCatalogs(ctx context.Context, direct, orm catalog) error {
	for _, variant := range []struct {
		name    string
		catalog catalog
	}{{"sql", direct}, {"gorm", orm}} {
		if err := verifyCatalog(ctx, variant.catalog); err != nil {
			return fmt.Errorf("%s preflight: %w", variant.name, err)
		}
	}

	directPage, err := direct.Page(ctx, searchName, pageSize)
	if err != nil {
		return err
	}
	ormPage, err := orm.Page(ctx, searchName, pageSize)
	if err != nil {
		return err
	}
	if !slices.Equal(directPage, ormPage) {
		return errors.New("SQL and GORM returned different catalog pages")
	}
	return nil
}

func verifyCatalog(ctx context.Context, c catalog) error {
	want := newMedicine(1, "Демо % концентрат")
	want.ExpireTime = 13
	if err := verifyMedicine(ctx, c, want); err != nil {
		return err
	}

	for _, sample := range []struct {
		filter string
		count  int64
	}{{"аЛьФа", 2500}, {"%", 1}, {"_", 1}, {"", 10000}} {
		count, err := c.Count(ctx, sample.filter)
		if err != nil {
			return err
		}
		if count != sample.count {
			return fmt.Errorf("count for %q: got %d, want %d", sample.filter, count, sample.count)
		}
	}
	page, err := c.Page(ctx, searchName, pageSize)
	if err != nil {
		return err
	}
	if len(page) != pageSize || page[0].Name != "Демо Альфа 00084" {
		return errors.New("unexpected second page of the catalog")
	}

	want = newMedicine(10001, "Демо проверка")
	if err := c.Insert(ctx, want); err != nil {
		return err
	}
	if err := verifyMedicine(ctx, c, want); err != nil {
		return err
	}
	want.Name = "Демо проверка переименования"
	if err := c.Rename(ctx, want.ID, want.Name); err != nil {
		return err
	}
	if err := verifyMedicine(ctx, c, want); err != nil {
		return err
	}
	if err := c.Delete(ctx, want.ID); err != nil {
		return err
	}
	if _, err := c.Find(ctx, want.ID); !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("expected no rows after deletion, got %v", err)
	}
	return nil
}

func verifyMedicine(ctx context.Context, c catalog, want entity.MedicineEntity) error {
	got, err := c.Find(ctx, want.ID)
	if err != nil {
		return err
	}
	if got != want {
		return fmt.Errorf("medicine %s: got %+v, want %+v", want.ID, got, want)
	}
	return nil
}

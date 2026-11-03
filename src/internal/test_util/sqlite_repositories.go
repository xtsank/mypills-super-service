package testutil

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	pgrepo "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/repository"
)

func openSQLite(t *testing.T, schema []string) *sqlx.DB {
	t.Helper()
	db, err := sqlx.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.Ping())
	_, err = db.Exec(sqliteEnableForeignKeys)
	require.NoError(t, err)
	for _, query := range schema {
		_, err = db.Exec(query)
		require.NoError(t, err)
	}
	return db
}

func NewMedicineRepository(t *testing.T) medicine.IMedicineRepository {
	t.Helper()
	injector := do.New()
	do.ProvideValue(injector, openSQLite(t, sqliteMedicineSchema))
	repo, err := pgrepo.NewPostgresMedicineRepository(injector)
	require.NoError(t, err)
	return repo
}

func NewCabinetRepository(t *testing.T) cabinet_item.ICabinetItemRepository {
	t.Helper()
	injector := do.New()
	do.ProvideValue(injector, openSQLite(t, sqliteCabinetSchema))
	repo, err := pgrepo.NewPostgresLegacyCabinetItemRepository(injector)
	require.NoError(t, err)
	return repo
}

func SeedMedicines(t *testing.T, repo medicine.IMedicineRepository, medicines ...*medicine.Medicine) {
	t.Helper()
	for _, med := range medicines {
		require.NoError(t, repo.Create(context.Background(), med))
	}
}

func SeedCabinetItems(t *testing.T, repo cabinet_item.ICabinetItemRepository, items ...*cabinet_item.CabinetItem) {
	t.Helper()
	for _, item := range items {
		require.NoError(t, repo.Save(context.Background(), item))
	}
}

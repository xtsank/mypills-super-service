//go:build integration

package integration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	"github.com/xtsank/mypills-super-service/src/internal/domain/restock"
	pgrepo "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/repository"
	"github.com/xtsank/mypills-super-service/src/internal/service"
)

// IDs of rows created by tests/fixtures/integration_seed.sql.
var (
	seedCabinetCreateUser = uuid.MustParse("00000000-0000-4000-8000-000000000001")
	seedCabinetOwner      = uuid.MustParse("00000000-0000-4000-8000-000000000002")
	seedCabinetOther      = uuid.MustParse("00000000-0000-4000-8000-000000000003")
	seedRestockUser       = uuid.MustParse("00000000-0000-4000-8000-000000000004")
	seedCabinet           = uuid.MustParse("00000000-0000-4000-8000-000000000011")
	seedMedicine          = uuid.MustParse("00000000-0000-4000-8000-000000000031")
)

func integrationDB(t *testing.T) *sqlx.DB {
	t.Helper()
	require.Equal(t, "mypills_test", os.Getenv("DB_NAME"), "integration tests must use the test database")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	db, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}

func integrationCabinetService(t *testing.T, db *sqlx.DB) (service.ICabinetService, cabinet.ICabinetRepository) {
	t.Helper()
	injector := do.New()
	do.ProvideValue(injector, db)
	do.Provide(injector, pgrepo.NewPostgresCabinetRepository)
	repo := do.MustInvoke[cabinet.ICabinetRepository](injector)
	svc, err := service.NewCabinetService(injector)
	require.NoError(t, err)
	return svc, repo
}

func integrationRestockService(t *testing.T, db *sqlx.DB) (service.IRestockService, restock.IRestockRepository) {
	t.Helper()
	injector := do.New()
	do.ProvideValue(injector, db)
	do.Provide(injector, pgrepo.NewPostgresMedicineRepository)
	do.Provide(injector, pgrepo.NewPostgresRestockRepository)
	repo := do.MustInvoke[restock.IRestockRepository](injector)
	svc, err := service.NewRestockService(injector)
	require.NoError(t, err)
	return svc, repo
}

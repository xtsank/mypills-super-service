package helpers

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
	pgrepo "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/repository"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
)

type ClassicRepositories struct {
	DB       *sqlx.DB
	SQL      sqlmock.Sqlmock
	User     user.IUserRepository
	Medicine medicine.IMedicineRepository
	Cabinet  cabinet_item.ICabinetItemRepository
}

func NewClassicRepositories(t *testing.T) ClassicRepositories {
	t.Helper()

	rawDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	db := sqlx.NewDb(rawDB, "sqlmock")
	t.Cleanup(func() {
		sqlMock.ExpectClose()
		require.NoError(t, db.Close())
	})

	injector := do.New()
	do.ProvideValue(injector, db)

	userRepo, err := pgrepo.NewPostgresUserRepository(injector)
	require.NoError(t, err)
	medicineRepo, err := pgrepo.NewPostgresMedicineRepository(injector)
	require.NoError(t, err)
	cabinetRepo, err := pgrepo.NewPostgresCabinetItemRepository(injector)
	require.NoError(t, err)

	return ClassicRepositories{
		DB:       db,
		SQL:      sqlMock,
		User:     userRepo,
		Medicine: medicineRepo,
		Cabinet:  cabinetRepo,
	}
}

func ExpectCabinetItemLookup(sqlMock sqlmock.Sqlmock, userID, medicineID uuid.UUID, date time.Time, item *cabinet_item.CabinetItem) {
	expectation := sqlMock.ExpectQuery(regexp.QuoteMeta(queries.CabinetItem.SelectExistingByKey)).
		WithArgs(userID, medicineID, date)
	if item == nil {
		expectation.WillReturnRows(testutil.CabinetRows())
		return
	}
	expectation.WillReturnRows(testutil.CabinetRows(item))
}

func ExpectCabinetUpdate(t *testing.T, repos ClassicRepositories, item *cabinet_item.CabinetItem) {
	t.Helper()

	query, _, err := sqlx.Named(queries.CabinetItem.UpdateQuantity, map[string]any{
		"id":       item.ID,
		"quantity": item.Quantity,
	})
	require.NoError(t, err)
	query = repos.DB.Rebind(query)

	repos.SQL.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(item.Quantity, item.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func ExpectCabinetInsert(t *testing.T, repos ClassicRepositories, item *cabinet_item.CabinetItem) {
	t.Helper()

	query, _, err := sqlx.Named(queries.CabinetItem.InsertItem, map[string]any{
		"id":                  uuid.Nil,
		"user_id":             item.UserID,
		"medicine_id":         item.MedicineID,
		"date_of_manufacture": item.DateOfManufacture,
		"quantity":            item.Quantity,
	})
	require.NoError(t, err)
	query = repos.DB.Rebind(query)

	repos.SQL.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(sqlmock.AnyArg(), item.UserID, item.MedicineID, item.DateOfManufacture, item.Quantity).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func ExpectUserByID(sqlMock sqlmock.Sqlmock, value *user.User) {
	sqlMock.ExpectQuery(regexp.QuoteMeta(queries.User.FindByID)).
		WithArgs(value.ID).
		WillReturnRows(testutil.UserBaseRows(value))
	sqlMock.ExpectQuery(regexp.QuoteMeta(queries.User.SelectIllnesses)).
		WithArgs(value.ID).
		WillReturnRows(testutil.UUIDRows("illness_id", value.Illnesses...))
	sqlMock.ExpectQuery(regexp.QuoteMeta(queries.User.SelectAllergies)).
		WithArgs(value.ID).
		WillReturnRows(testutil.UUIDRows("substance_id", value.Allergies...))
}

func ExpectMedicinesByIllness(sqlMock sqlmock.Sqlmock, illnessID uuid.UUID, values []*medicine.Medicine) {
	sqlMock.ExpectQuery(regexp.QuoteMeta(queries.Medicine.SelectByIllness)).
		WithArgs(illnessID).
		WillReturnRows(testutil.MedicineBaseRows(values...))

	for _, value := range values {
		sqlMock.ExpectQuery(regexp.QuoteMeta(queries.Medicine.SelectSubstances)).
			WithArgs(value.ID).
			WillReturnRows(testutil.SubstanceRows(value))
		sqlMock.ExpectQuery(regexp.QuoteMeta(queries.Medicine.SelectDosages)).
			WithArgs(value.ID).
			WillReturnRows(testutil.DosageRows(value))
		sqlMock.ExpectQuery(regexp.QuoteMeta(queries.Medicine.SelectContraindications)).
			WithArgs(value.ID).
			WillReturnRows(testutil.UUIDRows("illness_id", value.Contraindications...))
		sqlMock.ExpectQuery(regexp.QuoteMeta(queries.Medicine.SelectRecommendations)).
			WithArgs(value.ID).
			WillReturnRows(testutil.UUIDRows("illness_id", value.Recommendation...))
	}
}

func ExpectCabinetByUserID(sqlMock sqlmock.Sqlmock, userID uuid.UUID, items []*cabinet_item.CabinetItem) {
	sqlMock.ExpectQuery(regexp.QuoteMeta(queries.CabinetItem.SelectByUserID)).
		WithArgs(userID).
		WillReturnRows(testutil.CabinetRows(items...))
}

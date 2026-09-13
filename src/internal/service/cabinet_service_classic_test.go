package service

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
	testhelpers "github.com/xtsank/mypills-super-service/src/internal/test_util/helpers"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/testid"
)

func TestClassic_CabinetService_AddItem_UpdatesExisting(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			repos := testhelpers.NewClassicRepositories(t)
			svc := &CabinetService{cabinetRepo: repos.Cabinet}
			userID := testid.User
			medicineID := testid.Safe
			manufacturedAt := time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC)
			existing := testutil.NewCabinetItemBuilder().
				WithID(testid.CabinetSafe).
				WithUserID(userID).
				WithMedicineID(medicineID).
				WithManufactureDate(manufacturedAt).
				WithQuantity(5).
				Build()
			expectedUpdate := *existing
			expectedUpdate.Quantity = 8
			testhelpers.ExpectCabinetItemLookup(repos.SQL, userID, medicineID, manufacturedAt, existing)
			testhelpers.ExpectCabinetUpdate(t, repos, &expectedUpdate)

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.AddItem(ctx, command.NewAddItemCmd(userID, medicineID, manufacturedAt, 3))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.NoError(t, err)
					require.NotNil(t, result)
					require.Equal(t, existing.ID, result.ID)
					require.Equal(t, float32(8), result.Quantity)
					require.NoError(t, repos.SQL.ExpectationsWereMet())
				})
			})
		})
	}, allure.WithLabel("style", "classic"), allure.WithLabel("technique", "state transition"), allure.WithLabel("service", "CabinetService"))
}

func TestClassic_CabinetService_AddItem_CreatesNew(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			repos := testhelpers.NewClassicRepositories(t)
			svc := &CabinetService{cabinetRepo: repos.Cabinet}
			userID := testid.User
			medicineID := testid.Safe
			manufacturedAt := time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC)
			expected := testutil.NewCabinetItemBuilder().
				WithUserID(userID).
				WithMedicineID(medicineID).
				WithManufactureDate(manufacturedAt).
				WithQuantity(2).
				Build()
			testhelpers.ExpectCabinetItemLookup(repos.SQL, userID, medicineID, manufacturedAt, nil)
			testhelpers.ExpectCabinetInsert(t, repos, expected)

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.AddItem(ctx, command.NewAddItemCmd(userID, medicineID, manufacturedAt, 2))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.NoError(t, err)
					require.NotNil(t, result)
					require.NotEqual(t, uuid.Nil, result.ID)
					require.Equal(t, medicineID, result.MedicineID)
					require.Equal(t, float32(2), result.Quantity)
					require.NoError(t, repos.SQL.ExpectationsWereMet())
				})
			})
		})
	}, allure.WithLabel("style", "classic"), allure.WithLabel("technique", "equivalence partitioning; decision table"), allure.WithLabel("service", "CabinetService"))
}

func TestClassic_CabinetService_AddItem_RejectsZeroQuantity(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			repos := testhelpers.NewClassicRepositories(t)
			svc := &CabinetService{cabinetRepo: repos.Cabinet}
			userID := testid.User
			medicineID := testid.Safe
			manufacturedAt := time.Date(2025, time.March, 5, 0, 0, 0, 0, time.UTC)
			testhelpers.ExpectCabinetItemLookup(repos.SQL, userID, medicineID, manufacturedAt, nil)

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.AddItem(ctx, command.NewAddItemCmd(userID, medicineID, manufacturedAt, 0))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.Nil(t, result)
					require.Error(t, err)
					require.True(t, stderrors.Is(err, apperrors.ErrQtyTooLow))
					require.NoError(t, repos.SQL.ExpectationsWereMet())
				})
			})
		})
	}, allure.WithLabel("style", "classic"), allure.WithLabel("technique", "boundary value; negative"), allure.WithLabel("service", "CabinetService"))
}

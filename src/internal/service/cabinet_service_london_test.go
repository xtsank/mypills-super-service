package service

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/mocks"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/testid"
)

func TestLondon_CabinetService_AddItem_UpdatesExisting(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			repo := new(mocks.MockCabinetRepository)
			svc := &CabinetService{cabinetRepo: repo}
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

			repo.On("FindExistingCabinetItem", ctx, userID, medicineID, manufacturedAt).
				Return(existing, nil).
				Once()
			repo.On("Update", ctx, mock.MatchedBy(func(value *cabinet_item.CabinetItem) bool {
				return value == existing && value.Quantity == 8
			})).Return(nil).Once()

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.AddItem(ctx, command.NewAddItemCmd(userID, medicineID, manufacturedAt, 3))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.NoError(t, err)
					require.NotNil(t, result)
					require.Equal(t, existing.ID, result.ID)
					require.Equal(t, medicineID, result.MedicineID)
					require.Equal(t, float32(8), result.Quantity)
					repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
					repo.AssertExpectations(t)
				})
			})
		})
	}, allure.WithLabel("style", "london"), allure.WithLabel("technique", "state transition"), allure.WithLabel("service", "CabinetService"))
}

func TestLondon_CabinetService_AddItem_CreatesNew(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			repo := new(mocks.MockCabinetRepository)
			svc := &CabinetService{cabinetRepo: repo}
			userID := testid.User
			medicineID := testid.Safe
			manufacturedAt := time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC)

			repo.On("FindExistingCabinetItem", ctx, userID, medicineID, manufacturedAt).
				Return(nil, nil).
				Once()
			repo.On("Save", ctx, mock.MatchedBy(func(value *cabinet_item.CabinetItem) bool {
				return value.ID != uuid.Nil && value.UserID == userID && value.MedicineID == medicineID &&
					value.DateOfManufacture.Equal(manufacturedAt) && value.Quantity == 2
			})).Return(nil).Once()

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
					repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
					repo.AssertExpectations(t)
				})
			})
		})
	}, allure.WithLabel("style", "london"), allure.WithLabel("technique", "equivalence partitioning; decision table"), allure.WithLabel("service", "CabinetService"))
}

func TestLondon_CabinetService_AddItem_RejectsZeroQuantity(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			repo := new(mocks.MockCabinetRepository)
			svc := &CabinetService{cabinetRepo: repo}
			userID := testid.User
			medicineID := testid.Safe
			manufacturedAt := time.Date(2025, time.March, 5, 0, 0, 0, 0, time.UTC)

			repo.On("FindExistingCabinetItem", ctx, userID, medicineID, manufacturedAt).
				Return(nil, nil).
				Once()

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.AddItem(ctx, command.NewAddItemCmd(userID, medicineID, manufacturedAt, 0))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.Nil(t, result)
					require.Error(t, err)
					require.True(t, stderrors.Is(err, apperrors.ErrQtyTooLow))
					repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
					repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
					repo.AssertExpectations(t)
				})
			})
		})
	}, allure.WithLabel("style", "london"), allure.WithLabel("technique", "boundary value; negative"), allure.WithLabel("service", "CabinetService"))
}

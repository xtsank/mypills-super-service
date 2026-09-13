package service

import (
	"context"
	stderrors "errors"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/fixtures"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/mocks"
)

func TestLondon_MedicineService_Select_FiltersAndCalculatesDosage(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			fixture := fixtures.NewMedicineSelection()
			userRepo := new(mocks.MockUserRepository)
			medicineRepo := new(mocks.MockMedicineRepository)
			cabinetRepo := new(mocks.MockCabinetRepository)
			svc := &MedicineService{userRepo: userRepo, medicineRepo: medicineRepo, cabinetRepo: cabinetRepo}

			userRepo.On("FindByID", ctx, fixture.UserID).Return(fixture.User, nil).Once()
			medicineRepo.On("FindByIllness", ctx, fixture.IllnessID).Return(fixture.Medicines, nil).Once()
			cabinetRepo.On("FindByUserID", ctx, fixture.UserID).Return(fixture.Cabinet, nil).Once()

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.Select(ctx, command.NewSelectMedicineCmd(fixture.UserID, fixture.IllnessID))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.NoError(t, err)
					require.NotNil(t, result)
					require.Len(t, result.Recommendations, 1)
					recommendation := result.Recommendations[0]
					require.Equal(t, fixture.ExpectedMed.ID, recommendation.ID)
					require.Equal(t, fixture.ExpectedMed.Name, recommendation.Name)
					require.Equal(t, fixture.ExpectedMed.MethodOfApplication, recommendation.MethodOfApplication)
					require.Equal(t, float32(2.5), recommendation.Dosage)
					require.Equal(t, 3, recommendation.Frequency)
					require.Equal(t, float32(4), recommendation.QuantityInCabinet)
					userRepo.AssertExpectations(t)
					medicineRepo.AssertExpectations(t)
					cabinetRepo.AssertExpectations(t)
				})
			})
		})
	}, allure.WithLabel("style", "london"), allure.WithLabel("technique", "combinatorial; decision table; boundary value"), allure.WithLabel("service", "MedicineService"))
}

func TestLondon_MedicineService_Select_PropagatesRepositoryError(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			fixture := fixtures.NewMedicineSelection()
			userRepo := new(mocks.MockUserRepository)
			medicineRepo := new(mocks.MockMedicineRepository)
			cabinetRepo := new(mocks.MockCabinetRepository)
			svc := &MedicineService{userRepo: userRepo, medicineRepo: medicineRepo, cabinetRepo: cabinetRepo}
			databaseErr := stderrors.New("medicine query failed")
			repositoryErr := apperrors.ErrInternal.WithError(databaseErr)

			userRepo.On("FindByID", ctx, fixture.UserID).Return(fixture.User, nil).Once()
			medicineRepo.On("FindByIllness", ctx, fixture.IllnessID).Return(nil, repositoryErr).Once()

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.Select(ctx, command.NewSelectMedicineCmd(fixture.UserID, fixture.IllnessID))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.Nil(t, result)
					require.Error(t, err)
					require.True(t, stderrors.Is(err, apperrors.ErrInternal))
					require.True(t, stderrors.Is(err, databaseErr))
					cabinetRepo.AssertNotCalled(t, "FindByUserID", mock.Anything, mock.Anything)
					userRepo.AssertExpectations(t)
					medicineRepo.AssertExpectations(t)
				})
			})
		})
	}, allure.WithLabel("style", "london"), allure.WithLabel("technique", "fault injection; negative"), allure.WithLabel("service", "MedicineService"))
}

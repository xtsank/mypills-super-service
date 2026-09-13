package service

import (
	"context"
	stderrors "errors"
	"regexp"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/stretchr/testify/require"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/queries"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/fixtures"
	testhelpers "github.com/xtsank/mypills-super-service/src/internal/test_util/helpers"
)

func TestClassic_MedicineService_Select_FiltersAndCalculatesDosage(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			fixture := fixtures.NewMedicineSelection()
			repos := testhelpers.NewClassicRepositories(t)
			svc := &MedicineService{userRepo: repos.User, medicineRepo: repos.Medicine, cabinetRepo: repos.Cabinet}
			testhelpers.ExpectUserByID(repos.SQL, fixture.User)
			testhelpers.ExpectMedicinesByIllness(repos.SQL, fixture.IllnessID, fixture.Medicines)
			testhelpers.ExpectCabinetByUserID(repos.SQL, fixture.UserID, fixture.Cabinet)

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
					require.NoError(t, repos.SQL.ExpectationsWereMet())
				})
			})
		})
	}, allure.WithLabel("style", "classic"), allure.WithLabel("technique", "combinatorial; decision table; boundary value"), allure.WithLabel("service", "MedicineService"))
}

func TestClassic_MedicineService_Select_PropagatesRepositoryError(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.Step("Arrange", func(_ *allure.Context) {
			// Arrange
			ctx := context.Background()
			fixture := fixtures.NewMedicineSelection()
			repos := testhelpers.NewClassicRepositories(t)
			svc := &MedicineService{userRepo: repos.User, medicineRepo: repos.Medicine, cabinetRepo: repos.Cabinet}
			databaseErr := stderrors.New("medicine query failed")
			testhelpers.ExpectUserByID(repos.SQL, fixture.User)
			repos.SQL.ExpectQuery(regexp.QuoteMeta(queries.Medicine.SelectByIllness)).
				WithArgs(fixture.IllnessID).
				WillReturnError(databaseErr)

			a.Step("Act", func(_ *allure.Context) {
				// Act
				result, err := svc.Select(ctx, command.NewSelectMedicineCmd(fixture.UserID, fixture.IllnessID))

				a.Step("Assert", func(_ *allure.Context) {
					// Assert
					require.Nil(t, result)
					require.Error(t, err)
					require.True(t, stderrors.Is(err, apperrors.ErrInternal))
					require.True(t, stderrors.Is(err, databaseErr))
					require.NoError(t, repos.SQL.ExpectationsWereMet())
				})
			})
		})
	}, allure.WithLabel("style", "classic"), allure.WithLabel("technique", "fault injection; negative"), allure.WithLabel("service", "MedicineService"))
}

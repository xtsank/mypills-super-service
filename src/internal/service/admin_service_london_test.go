package service

import (
	"context"
	stderrors "errors"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/mocks"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

func TestLondon_AdminService_UpdateMedicine_StopsOnReadError(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Ошибка поиска: без обновления")
		a.Description("Если поиск лекарства завершился ошибкой, сервис возвращает её и не вызывает обновление репозитория.")
		a.Parameter("Сбой в", "FindByID")
		var ctx context.Context
		var repo *mocks.MockMedicineRepository
		var svc *AdminService
		var lookupErr error
		var cmd *command.UpdateMedicineCmd
		a.Step("Arrange", func(a *allure.Context) {
			ctx = context.Background()
			repo = new(mocks.MockMedicineRepository)
			svc = &AdminService{medicineRepo: repo}
			lookupErr = stderrors.New("medicine lookup failed")
			repo.On("FindByID", ctx, testutil.MedicineA).Return(nil, lookupErr).Once()
			newExpiry := 36
			cmd = &command.UpdateMedicineCmd{ID: testutil.MedicineA, ExpireTime: &newExpiry}
		})
		var result *res.AdminResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.UpdateMedicine(ctx, cmd)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.Nil(t, result)
			require.ErrorIs(t, err, lookupErr)
			repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			repo.AssertExpectations(t)
		})
	}, allure.WithLabel("стиль", "лондонский"), allure.WithLabel("сервис", "AdminService"), allure.WithLabel("техника", "внедрение ошибки"), allure.WithLabel("подготовка", "mock; ID"), allure.WithEpic("MyPills"), allure.WithFeature("AdminService"), allure.WithStory("UpdateMedicine"), allure.WithParentSuite("MyPills"), allure.WithSuite("AdminService"), allure.WithSubSuite("лондонский"))
}

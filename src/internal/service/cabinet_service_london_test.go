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
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/mocks"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

func TestLondon_CabinetService_AddItem_CreatesNew(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Новая позиция: Save без Update")
		a.Description("Класс эквивалентности: позиции ещё нет, количество положительно и равно 2. По таблице решений сервис выбирает Save, а не Update; mock проверяет вызовы.")
		a.Parameter("Количество", "2")
		a.Attachment("Данные", []byte("Мок репозитория testify; фиксированные UUID и дата производства"), "text/plain")
		var ctx context.Context
		var repo *mocks.MockCabinetRepository
		var svc *LegacyCabinetService
		var userID, medicineID uuid.UUID
		var madeAt time.Time
		a.Step("Arrange", func(_ *allure.Context) {
			ctx = context.Background()
			repo = new(mocks.MockCabinetRepository)
			svc = &LegacyCabinetService{cabinetRepo: repo}
			userID = testutil.CabinetUserID
			medicineID = testutil.CabinetMedicineID
			madeAt = time.Date(2025, time.February, 10, 0, 0, 0, 0, time.UTC)
			repo.On("FindExistingCabinetItem", ctx, userID, medicineID, madeAt).Return(nil, nil).Once()
			repo.On("Save", ctx, mock.MatchedBy(func(item *cabinet_item.CabinetItem) bool {
				return item != nil && item.ID != uuid.Nil && item.UserID == userID &&
					item.MedicineID == medicineID && item.DateOfManufacture.Equal(madeAt) && item.Quantity == 2
			})).Return(nil).Once()
		})
		var result *res.CabinetResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.AddItem(ctx, command.NewAddItemCmd(userID, medicineID, madeAt, 2))
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotEqual(t, uuid.Nil, result.ID)
			require.Equal(t, medicineID, result.MedicineID)
			require.Equal(t, float32(2), result.Quantity)
			repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			repo.AssertExpectations(t)
		})
	}, allure.WithLabel("стиль", "лондонский"), allure.WithLabel("техника", "таблица решений"), allure.WithLabel("сервис", "CabinetService"), allure.WithLabel("подготовка", "mock; ID и дата"), allure.WithEpic("MyPills"), allure.WithFeature("CabinetService"), allure.WithStory("AddItem"), allure.WithParentSuite("MyPills"), allure.WithSuite("CabinetService"), allure.WithSubSuite("лондонский"))
}

func TestLondon_CabinetService_AddItem_StopsOnLookupError(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Ошибка поиска: без записи")
		a.Description("При ошибке начального поиска сервис возвращает ту же ошибку и ничего не записывает.")
		a.Parameter("Сбой в", "FindExistingCabinetItem")
		var ctx context.Context
		var repo *mocks.MockCabinetRepository
		var svc *LegacyCabinetService
		var userID, medicineID uuid.UUID
		var madeAt time.Time
		var lookupErr error
		a.Step("Arrange", func(_ *allure.Context) {
			ctx = context.Background()
			repo = new(mocks.MockCabinetRepository)
			svc = &LegacyCabinetService{cabinetRepo: repo}
			userID = testutil.CabinetUserID
			medicineID = testutil.CabinetMedicineID
			madeAt = time.Date(2025, time.April, 12, 0, 0, 0, 0, time.UTC)
			lookupErr = stderrors.New("cabinet lookup failed")
			repo.On("FindExistingCabinetItem", ctx, userID, medicineID, madeAt).Return(nil, lookupErr).Once()
		})
		var result *res.CabinetResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.AddItem(ctx, command.NewAddItemCmd(userID, medicineID, madeAt, 3))
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.Nil(t, result)
			require.ErrorIs(t, err, lookupErr)
			repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
			repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			repo.AssertExpectations(t)
		})
	}, allure.WithLabel("стиль", "лондонский"), allure.WithLabel("техника", "внедрение ошибки"), allure.WithLabel("сервис", "CabinetService"), allure.WithLabel("подготовка", "mock; ID и дата"), allure.WithEpic("MyPills"), allure.WithFeature("CabinetService"), allure.WithStory("AddItem"), allure.WithParentSuite("MyPills"), allure.WithSuite("CabinetService"), allure.WithSubSuite("лондонский"))
}

func TestLondon_CabinetService_RemoveItem_PropagatesDeleteError(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Ошибка Delete возвращается")
		a.Description("Ошибка Delete видна вызывающему коду; репозиторий вызывается один раз.")
		a.Parameter("Сбой в", "Delete")
		var ctx context.Context
		var repo *mocks.MockCabinetRepository
		var svc *LegacyCabinetService
		var itemID uuid.UUID
		var deleteErr error
		a.Step("Arrange", func(_ *allure.Context) {
			ctx = context.Background()
			repo = new(mocks.MockCabinetRepository)
			svc = &LegacyCabinetService{cabinetRepo: repo}
			itemID = testutil.CabinetItemID
			deleteErr = stderrors.New("cabinet delete failed")
			repo.On("Delete", ctx, itemID).Return(deleteErr).Once()
		})
		var err error
		a.Step("Act", func(_ *allure.Context) {
			err = svc.RemoveItem(ctx, command.NewRemoveItemCmd(itemID))
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.ErrorIs(t, err, deleteErr)
			repo.AssertExpectations(t)
		})
	}, allure.WithLabel("стиль", "лондонский"), allure.WithLabel("техника", "внедрение ошибки"), allure.WithLabel("сервис", "CabinetService"), allure.WithLabel("подготовка", "mock; ID"), allure.WithEpic("MyPills"), allure.WithFeature("CabinetService"), allure.WithStory("RemoveItem"), allure.WithParentSuite("MyPills"), allure.WithSuite("CabinetService"), allure.WithSubSuite("лондонский"))
}

func TestLondon_CabinetService_ListItems_PropagatesReadError(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Ошибка чтения списка возвращается")
		a.Description("При ошибке чтения списка сервис возвращает ошибку вместо списка DTO.")
		a.Parameter("Сбой в", "FindDetailsByUserID")
		var ctx context.Context
		var repo *mocks.MockCabinetRepository
		var svc *LegacyCabinetService
		var userID uuid.UUID
		var readErr error
		a.Step("Arrange", func(_ *allure.Context) {
			ctx = context.Background()
			repo = new(mocks.MockCabinetRepository)
			svc = &LegacyCabinetService{cabinetRepo: repo}
			userID = testutil.CabinetUserID
			readErr = stderrors.New("cabinet list failed")
			repo.On("FindDetailsByUserID", ctx, userID).Return(nil, readErr).Once()
		})
		var result []res.CabinetItemDetailsResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.ListItems(ctx, userID)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.Nil(t, result)
			require.ErrorIs(t, err, readErr)
			repo.AssertExpectations(t)
		})
	}, allure.WithLabel("стиль", "лондонский"), allure.WithLabel("техника", "внедрение ошибки"), allure.WithLabel("сервис", "CabinetService"), allure.WithLabel("подготовка", "mock; ID"), allure.WithEpic("MyPills"), allure.WithFeature("CabinetService"), allure.WithStory("ListItems"), allure.WithParentSuite("MyPills"), allure.WithSuite("CabinetService"), allure.WithSubSuite("лондонский"))
}

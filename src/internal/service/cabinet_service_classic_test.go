package service

import (
	"context"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

func TestClassic_CabinetService_UpdateQty_PersistsNewQuantity(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Количество аптечки: 5 → 3")
		a.Description("Изменение количества с 5 до 3 отражается и в возвращённом DTO, и в записи БД через настоящий репозиторий.")
		a.Parameter("Было", "5")
		a.Parameter("Стало", "3")
		var ctx context.Context
		var repo cabinet_item.ICabinetItemRepository
		var original *cabinet_item.CabinetItem
		var svc *LegacyCabinetService
		var cmd *command.UpdateQtyCmd
		a.Step("Arrange", func(a *allure.Context) {
			ctx = context.Background()
			repo = testutil.NewCabinetRepository(t)
			original = testutil.CabinetItemWithQuantityFive().Build()
			testutil.SeedCabinetItems(t, repo, original)
			svc = &LegacyCabinetService{cabinetRepo: repo}
			cmd = command.NewUpdateQtyCmd(original.ID, 3)
		})
		var result *res.CabinetResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.UpdateQty(ctx, cmd)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, original.ID, result.ID)
			require.Equal(t, original.MedicineID, result.MedicineID)
			require.Equal(t, float32(3), result.Quantity)

			stored, readErr := repo.FindById(ctx, original.ID)
			require.NoError(t, readErr)
			require.NotNil(t, stored)
			require.Equal(t, float32(3), stored.Quantity)
			require.Equal(t, original.DateOfManufacture, stored.DateOfManufacture)
		})
	}, allure.WithLabel("стиль", "классический"), allure.WithLabel("сервис", "CabinetService"), allure.WithLabel("техника", "переход состояния"), allure.WithLabel("подготовка", "Mother; SQLite"), allure.WithEpic("MyPills"), allure.WithFeature("CabinetService"), allure.WithStory("UpdateQty"), allure.WithParentSuite("MyPills"), allure.WithSuite("CabinetService"), allure.WithSubSuite("классический"))
}

func TestClassic_CabinetService_UpdateQty_MissingItemPreservesGroup(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Ошибка обновления отсутствующей позиции")
		a.Description("Обновление неизвестного ID возвращает ошибку; две существующие позиции, подготовленные одной группой, остаются без изменений.")
		var repo cabinet_item.ICabinetItemRepository
		var svc *LegacyCabinetService
		var first, second *cabinet_item.CabinetItem
		var cmd *command.UpdateQtyCmd
		a.Step("Arrange", func(_ *allure.Context) {
			repo = testutil.NewCabinetRepository(t)
			group := testutil.CabinetItemsForMissingUpdate()
			first = group.Existing.Build()
			second = group.Other.Build()
			testutil.SeedCabinetItems(t, repo, first, second)
			svc = &LegacyCabinetService{cabinetRepo: repo}
			cmd = command.NewUpdateQtyCmd(group.MissingID, 3)
		})
		var result *res.CabinetResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.UpdateQty(context.Background(), cmd)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.Nil(t, result)
			require.ErrorIs(t, err, apperrors.ErrCabinetItemNotFound)
			for _, expected := range []*cabinet_item.CabinetItem{first, second} {
				stored, readErr := repo.FindById(context.Background(), expected.ID)
				require.NoError(t, readErr)
				require.NotNil(t, stored)
				require.Equal(t, expected, stored)
			}
		})
	}, allure.WithLabel("стиль", "классический"), allure.WithLabel("сервис", "CabinetService"), allure.WithLabel("техника", "ошибка отсутствующего ID"), allure.WithLabel("подготовка", "Mother; Builder; группа SQLite"), allure.WithEpic("MyPills"), allure.WithFeature("CabinetService"), allure.WithStory("UpdateQty"), allure.WithParentSuite("MyPills"), allure.WithSuite("CabinetService"), allure.WithSubSuite("классический"))
}

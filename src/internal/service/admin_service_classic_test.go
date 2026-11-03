package service

import (
	"context"
	"testing"

	allure "github.com/allure-framework/allure-go/commons/gotest"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

func TestClassic_AdminService_AddMedicine_PersistsMedicineAndRelations(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Добавление лекарства со связями")
		a.Description("Новый препарат с веществом, дозировкой, показанием и противопоказанием сохраняется через настоящий SQL-репозиторий.")
		a.Parameter("Название", "Complex remedy")
		a.Parameter("Концентрация", "12,5")
		var ctx context.Context
		var repo medicine.IMedicineRepository
		var svc *AdminService
		var cmd *command.AddMedicineCmd
		a.Step("Arrange", func(a *allure.Context) {
			ctx = context.Background()
			repo = testutil.NewMedicineRepository(t)
			svc = &AdminService{medicineRepo: repo}
			cmd = &command.AddMedicineCmd{
				Name: "Complex remedy", ExpireTime: 18, MethodOfApplication: "oral",
				Form: testutil.FormA, Unit: testutil.UnitA,
				Recommendation:    []uuid.UUID{testutil.IllnessA},
				Contraindications: []uuid.UUID{testutil.IllnessB},
				Substances:        []*req.ActiveSubstanceDto{{ID: testutil.SubstanceA, Concentration: 12.5}},
				Dosages:           []*req.DosageRuleDto{{ValueFrom: 70, ValueTo: 90, Type: req.ByWeight, DosageValue: 2.5, NumberOfDosesPerDay: 2}},
			}
		})
		var result *res.AdminResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.AddMedicine(ctx, cmd)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotEqual(t, uuid.Nil, result.ID)
			require.Equal(t, "Complex remedy", result.Name)
			require.Equal(t, []uuid.UUID{testutil.IllnessA}, result.Recommendations)
			require.Equal(t, []uuid.UUID{testutil.IllnessB}, result.Contraindications)
			require.Len(t, result.Substances, 1)
			require.Equal(t, float32(12.5), result.Substances[0].Concentration)
			require.Len(t, result.Dosages, 1)
			require.Equal(t, float32(2.5), result.Dosages[0].DosageValue)

			stored, readErr := repo.FindByID(ctx, result.ID)
			require.NoError(t, readErr)
			require.Equal(t, "Complex remedy", stored.Name)
			require.Equal(t, []uuid.UUID{testutil.IllnessA}, stored.Recommendation)
			require.Equal(t, []uuid.UUID{testutil.IllnessB}, stored.Contraindications)
			require.Len(t, stored.Substances, 1)
			require.Equal(t, testutil.SubstanceA, stored.Substances[0].ID)
			require.Equal(t, float32(12.5), stored.Substances[0].Concentration)
			require.Len(t, stored.Dosages, 1)
			require.Equal(t, 70, stored.Dosages[0].ValueFrom)
			require.Equal(t, 90, stored.Dosages[0].ValueTo)
			require.Equal(t, float32(2.5), stored.Dosages[0].DosageValue)
			require.Equal(t, 2, stored.Dosages[0].NumberOfDosesPerDay)
		})
	}, allure.WithLabel("стиль", "классический"), allure.WithLabel("сервис", "AdminService"), allure.WithLabel("техника", "сохранение связей"), allure.WithLabel("подготовка", "SQLite; команда"), allure.WithEpic("MyPills"), allure.WithFeature("AdminService"), allure.WithStory("AddMedicine"), allure.WithParentSuite("MyPills"), allure.WithSuite("AdminService"), allure.WithSubSuite("классический"))
}

func TestClassic_AdminService_UpdateMedicine_ChangesOnlySuppliedFields(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Частичное обновление лекарства")
		a.Description("Явные значения false и новый срок действия меняются, а имя и связанные записи, не переданные в команде, сохраняются.")
		a.Parameter("Срок до", "24 месяца")
		a.Parameter("Срок после", "36 месяцев")
		var ctx context.Context
		var repo medicine.IMedicineRepository
		var original *medicine.Medicine
		var svc *AdminService
		var cmd *command.UpdateMedicineCmd
		a.Step("Arrange", func(a *allure.Context) {
			ctx = context.Background()
			repo = testutil.NewMedicineRepository(t)
			original = testutil.MedicineForPartialUpdate().Build()
			testutil.SeedMedicines(t, repo, original)
			svc = &AdminService{medicineRepo: repo}
			newExpiry := 36
			falseValue := false
			newMethod := "topical"
			cmd = &command.UpdateMedicineCmd{
				ID: original.ID, ExpireTime: &newExpiry, IsPrescription: &falseValue,
				EffectOnDriver: &falseValue, MethodOfApplication: &newMethod,
			}
		})
		var result *res.AdminResDto
		var err error
		a.Step("Act", func(_ *allure.Context) {
			result, err = svc.UpdateMedicine(ctx, cmd)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, 36, result.ExpireTime)
			require.False(t, result.IsPrescription)
			require.False(t, result.EffectOnDriver)
			require.Equal(t, "topical", result.MethodOfApplication)
			require.Equal(t, "Existing remedy", result.Name)

			stored, readErr := repo.FindByID(ctx, original.ID)
			require.NoError(t, readErr)
			require.Equal(t, 36, stored.ExpireTime)
			require.False(t, stored.IsPrescription)
			require.False(t, stored.EffectOnDriver)
			require.Equal(t, "topical", stored.MethodOfApplication)
			require.Equal(t, "Existing remedy", stored.Name)
			require.Equal(t, original.Substances, stored.Substances)
			require.Equal(t, original.Dosages, stored.Dosages)
		})
	}, allure.WithLabel("стиль", "классический"), allure.WithLabel("сервис", "AdminService"), allure.WithLabel("техника", "частичное обновление"), allure.WithLabel("подготовка", "Mother; SQLite"), allure.WithEpic("MyPills"), allure.WithFeature("AdminService"), allure.WithStory("UpdateMedicine"), allure.WithParentSuite("MyPills"), allure.WithSuite("AdminService"), allure.WithSubSuite("классический"))
}

func TestClassic_AdminService_UpdateComposition_ReplacesOnlySubstances(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Замена состава лекарства")
		a.Description("Замена состава в настоящей БД удаляет старое вещество, но сохраняет дозировку и показание.")
		a.Parameter("До", "A: 10")
		a.Parameter("После", "B: 7,5")
		var ctx context.Context
		var repo medicine.IMedicineRepository
		var original *medicine.Medicine
		var svc *AdminService
		var cmd command.UpdateCompositionCmd
		a.Step("Arrange", func(a *allure.Context) {
			ctx = context.Background()
			repo = testutil.NewMedicineRepository(t)
			original = testutil.MedicineWithComposition().Build()
			testutil.SeedMedicines(t, repo, original)
			svc = &AdminService{medicineRepo: repo}
			cmd = command.UpdateCompositionCmd{
				MedicineID: original.ID,
				Substances: []*req.ActiveSubstanceDto{{ID: testutil.SubstanceB, Concentration: 7.5}},
			}
		})
		var err error
		a.Step("Act", func(_ *allure.Context) {
			err = svc.UpdateComposition(ctx, cmd)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			stored, readErr := repo.FindByID(ctx, original.ID)
			require.NoError(t, readErr)
			require.Len(t, stored.Substances, 1)
			require.Equal(t, testutil.SubstanceB, stored.Substances[0].ID)
			require.Equal(t, float32(7.5), stored.Substances[0].Concentration)
			require.Equal(t, original.Dosages, stored.Dosages)
			require.Equal(t, original.Recommendation, stored.Recommendation)
		})
	}, allure.WithLabel("стиль", "классический"), allure.WithLabel("сервис", "AdminService"), allure.WithLabel("техника", "замена состава"), allure.WithLabel("подготовка", "Mother; SQLite"), allure.WithEpic("MyPills"), allure.WithFeature("AdminService"), allure.WithStory("UpdateComposition"), allure.WithParentSuite("MyPills"), allure.WithSuite("AdminService"), allure.WithSubSuite("классический"))
}

func TestClassic_AdminService_RemoveMedicine_RemovesMedicineWithRelations(t *testing.T) {
	allure.Wrap(t, func(a *allure.Context) {
		a.DisplayName("Удаление лекарства со связями")
		a.Description("Удаление препарата с составом, дозировкой и показанием проходит при включённых внешних ключах; после удаления препарат нельзя прочитать.")
		a.Parameter("Связи", "состав, доза, показание")
		var ctx context.Context
		var repo medicine.IMedicineRepository
		var original *medicine.Medicine
		var svc *AdminService
		var cmd *command.RemoveMedicineCmd
		a.Step("Arrange", func(a *allure.Context) {
			ctx = context.Background()
			repo = testutil.NewMedicineRepository(t)
			original = testutil.MedicineWithComposition().Build()
			testutil.SeedMedicines(t, repo, original)
			svc = &AdminService{medicineRepo: repo}
			cmd = command.NewRemoveMedicineCmd(original.ID)
		})
		var err error
		a.Step("Act", func(_ *allure.Context) {
			err = svc.RemoveMedicine(ctx, cmd)
		})
		a.Step("Assert", func(_ *allure.Context) {
			require.NoError(t, err)
			stored, readErr := repo.FindByID(ctx, original.ID)
			require.Nil(t, stored)
			require.ErrorIs(t, readErr, apperrors.ErrMedicineNotFound)
		})
	}, allure.WithLabel("стиль", "классический"), allure.WithLabel("сервис", "AdminService"), allure.WithLabel("техника", "удаление"), allure.WithLabel("подготовка", "Mother; SQLite"), allure.WithEpic("MyPills"), allure.WithFeature("AdminService"), allure.WithStory("RemoveMedicine"), allure.WithParentSuite("MyPills"), allure.WithSuite("AdminService"), allure.WithSubSuite("классический"))
}

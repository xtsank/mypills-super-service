package fixtures

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	testutil "github.com/xtsank/mypills-super-service/src/internal/test_util"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/testid"
)

type MedicineSelection struct {
	UserID      uuid.UUID
	IllnessID   uuid.UUID
	AllergyID   uuid.UUID
	User        *user.User
	Medicines   []*medicine.Medicine
	Cabinet     []*cabinet_item.CabinetItem
	ExpectedMed *medicine.Medicine
}

func NewMedicineSelection() MedicineSelection {
	userID := testid.User
	illnessID := testid.Illness
	allergyID := testid.Allergen

	safe := testutil.SafeMedicine(testid.Safe, illnessID)
	pregnancyRisk := testutil.NewMedicineBuilder().
		WithID(testid.Pregnancy).
		WithName("Pregnancy risk").
		WithPregnancyRisk(true).
		WithRecommendations(illnessID).
		Build()
	driverRisk := testutil.NewMedicineBuilder().
		WithID(testid.Driver).
		WithName("Driver risk").
		WithDriverRisk(true).
		WithRecommendations(illnessID).
		Build()
	contraindicated := testutil.NewMedicineBuilder().
		WithID(testid.Contra).
		WithName("Contraindicated").
		WithContraindications(illnessID).
		WithRecommendations(illnessID).
		Build()
	allergic := testutil.NewMedicineBuilder().
		WithID(testid.Allergic).
		WithName("Contains allergen").
		WithSubstances(medicine.ActiveSubstance{ID: allergyID, Concentration: 10}).
		WithRecommendations(illnessID).
		Build()
	absent := testutil.NewMedicineBuilder().
		WithID(testid.Absent).
		WithName("Not in cabinet").
		WithRecommendations(illnessID).
		Build()
	expired := testutil.NewMedicineBuilder().
		WithID(testid.Expired).
		WithName("Expired").
		WithExpireTime(1).
		WithRecommendations(illnessID).
		Build()
	allFactorsFalse := testutil.NewMedicineBuilder().
		WithID(testid.AllFalse).
		WithName("Unsafe and absent").
		WithPregnancyRisk(true).
		WithDriverRisk(true).
		WithContraindications(illnessID).
		WithSubstances(medicine.ActiveSubstance{ID: allergyID, Concentration: 5}).
		WithRecommendations(illnessID).
		Build()

	freshDate := time.Now().UTC().AddDate(0, -1, 0)
	cabinet := []*cabinet_item.CabinetItem{
		testutil.AvailableCabinetItem(testid.CabinetSafe, userID, safe.ID),
		testutil.NewCabinetItemBuilder().WithID(testid.CabinetPregnancy).WithUserID(userID).WithMedicineID(pregnancyRisk.ID).WithManufactureDate(freshDate).WithQuantity(2).Build(),
		testutil.NewCabinetItemBuilder().WithID(testid.CabinetDriver).WithUserID(userID).WithMedicineID(driverRisk.ID).WithManufactureDate(freshDate).WithQuantity(2).Build(),
		testutil.NewCabinetItemBuilder().WithID(testid.CabinetContra).WithUserID(userID).WithMedicineID(contraindicated.ID).WithManufactureDate(freshDate).WithQuantity(2).Build(),
		testutil.NewCabinetItemBuilder().WithID(testid.CabinetAllergic).WithUserID(userID).WithMedicineID(allergic.ID).WithManufactureDate(freshDate).WithQuantity(2).Build(),
		testutil.NewCabinetItemBuilder().WithID(testid.CabinetExpired).WithUserID(userID).WithMedicineID(expired.ID).WithManufactureDate(time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)).WithQuantity(2).Build(),
	}

	return MedicineSelection{
		UserID:    userID,
		IllnessID: illnessID,
		AllergyID: allergyID,
		User:      testutil.TestUser(userID, illnessID, allergyID),
		Medicines: []*medicine.Medicine{
			safe,
			pregnancyRisk,
			driverRisk,
			contraindicated,
			allergic,
			absent,
			expired,
			allFactorsFalse,
		},
		Cabinet:     cabinet,
		ExpectedMed: safe,
	}
}

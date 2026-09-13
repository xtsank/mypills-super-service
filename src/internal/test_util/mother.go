package testutil

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/test_util/testid"
)

func TestUser(id, illnessID, allergyID uuid.UUID) *user.User {
	return NewUserBuilder().
		WithID(id).
		WithWeight(70).
		WithPregnancy(true).
		WithDriver(true).
		WithIllnesses(illnessID).
		WithAllergies(allergyID).
		Build()
}

func SafeMedicine(id, illnessID uuid.UUID) *medicine.Medicine {
	return NewMedicineBuilder().
		WithID(id).
		WithName("Safe medicine").
		WithRecommendations(illnessID).
		WithDosages(medicine.DosageRule{
			ID:                  testid.DefaultDosageRule,
			ValueFrom:           70,
			ValueTo:             90,
			Type:                medicine.ByWeight,
			DosageValue:         2.5,
			NumberOfDosesPerDay: 3,
		}).
		Build()
}

func AvailableCabinetItem(id, userID, medicineID uuid.UUID) *cabinet_item.CabinetItem {
	return NewCabinetItemBuilder().
		WithID(id).
		WithUserID(userID).
		WithMedicineID(medicineID).
		WithManufactureDate(time.Now().UTC().AddDate(0, -1, 0)).
		WithQuantity(4).
		Build()
}

package testutil

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
)

func MedicineWithComposition() *MedicineBuilder {
	return NewMedicineBuilder().
		WithSubstances([]medicine.ActiveSubstance{{ID: SubstanceA, Concentration: 10}}).
		WithDosages([]medicine.DosageRule{{
			ID: DosageA, ValueFrom: 70, ValueTo: 90, Type: medicine.ByWeight,
			DosageValue: 2, NumberOfDosesPerDay: 2,
		}}).
		WithRecommendations([]uuid.UUID{IllnessA})
}

func MedicineForPartialUpdate() *MedicineBuilder {
	return NewMedicineBuilder().
		WithID(MedicineB).
		WithName("Existing remedy").
		WithSubstances([]medicine.ActiveSubstance{{ID: SubstanceA, Concentration: 10}}).
		WithDosages([]medicine.DosageRule{{
			ID: DosageA, ValueFrom: 70, ValueTo: 90, Type: medicine.ByWeight,
			DosageValue: 2, NumberOfDosesPerDay: 2,
		}})
}

func CabinetItemWithQuantityFive() *CabinetItemBuilder {
	return NewCabinetItemBuilder().WithQuantity(5)
}

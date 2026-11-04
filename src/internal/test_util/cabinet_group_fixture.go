package testutil

import "github.com/google/uuid"

type CabinetItemsFixture struct {
	Existing  *CabinetItemBuilder
	Other     *CabinetItemBuilder
	MissingID uuid.UUID
}

func CabinetItemsForMissingUpdate() CabinetItemsFixture {
	return CabinetItemsFixture{
		Existing: CabinetItemWithQuantityFive(),
		Other: NewCabinetItemBuilder().
			WithID(CabinetOtherItemID).
			WithMedicineID(MedicineB).
			WithQuantity(7),
		MissingID: CabinetMissingID,
	}
}

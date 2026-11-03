package testutil

import "github.com/google/uuid"

var (
	CabinetUserID      = fromName("cabinet-user")
	CabinetMedicineID  = fromName("cabinet-medicine")
	CabinetItemID      = fromName("cabinet-item")
	CabinetOtherItemID = fromName("cabinet-other-item")
	CabinetMissingID   = fromName("cabinet-missing-item")

	FormA      = fromName("form-a")
	UnitA      = fromName("unit-a")
	IllnessA   = fromName("illness-a")
	IllnessB   = fromName("illness-b")
	SubstanceA = fromName("substance-a")
	SubstanceB = fromName("substance-b")
	MedicineA  = fromName("medicine-a")
	MedicineB  = fromName("medicine-b")
	DosageA    = fromName("dosage-a")
)

func fromName(name string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("mypills-tests/"+name))
}

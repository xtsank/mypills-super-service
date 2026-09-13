package testid

import "github.com/google/uuid"

var (
	DefaultUser        = fromName("default-user")
	DefaultMedicine    = fromName("default-medicine")
	DefaultForm        = fromName("default-form")
	DefaultUnit        = fromName("default-unit")
	DefaultCabinetItem = fromName("default-cabinet-item")
	DefaultDosageRule  = fromName("default-dosage-rule")

	User      = fromName("medicine-selection-user")
	Illness   = fromName("medicine-selection-illness")
	Allergen  = fromName("medicine-selection-allergen")
	Safe      = fromName("safe-medicine")
	Pregnancy = fromName("pregnancy-risk-medicine")
	Driver    = fromName("driver-risk-medicine")
	Contra    = fromName("contraindicated-medicine")
	Allergic  = fromName("allergic-medicine")
	Absent    = fromName("absent-medicine")
	Expired   = fromName("expired-medicine")
	AllFalse  = fromName("all-factors-false-medicine")

	CabinetSafe      = fromName("cabinet-safe")
	CabinetPregnancy = fromName("cabinet-pregnancy-risk")
	CabinetDriver    = fromName("cabinet-driver-risk")
	CabinetContra    = fromName("cabinet-contraindicated")
	CabinetAllergic  = fromName("cabinet-allergic")
	CabinetExpired   = fromName("cabinet-expired")
)

func fromName(name string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("mypills-tests/"+name))
}

package entity

import "github.com/google/uuid"

type MedicineEntity struct {
	ID         uuid.UUID `db:"id"`
	FormId     uuid.UUID `db:"form_id"`
	UnitId     uuid.UUID `db:"unit_id"`
	Name       string    `db:"name"`
	ExpireTime int       `db:"expire_time"`

	EffectOnDriver      bool   `db:"effect_on_driver"`
	EffectOnPregnant    bool   `db:"effect_on_pregnant"`
	MethodOfApplication string `db:"method_of_application"`
	IsPrescription      bool   `db:"is_prescription"`
}

type DosageEntity struct {
	ID                  uuid.UUID `db:"id"`
	MedicineId          uuid.UUID `db:"medicine_id"`
	ValueFrom           int       `db:"value_from"`
	ValueTo             int       `db:"value_to"`
	DosageType          string    `db:"dosage_type"`
	DosageValue         float32   `db:"dosage_value"`
	NumberOfDosesPerDay int       `db:"number_of_doses_per_day"`
}

type MedicineIllnessEntity struct {
	ID         uuid.UUID `db:"id"`
	MedicineId uuid.UUID `db:"medicine_id"`
	IllnessId  uuid.UUID `db:"illness_id"`
}

type MedicineSubstanceEntity struct {
	ID            uuid.UUID `db:"id"`
	MedicineId    uuid.UUID `db:"medicine_id"`
	SubstanceId   uuid.UUID `db:"substance_id"`
	Concentration float32   `db:"concentration"`
}

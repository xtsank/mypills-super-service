package entity

import "github.com/google/uuid"

type MedicineEntity struct {
	ID                  uuid.UUID   `bson:"_id"`
	FormID              uuid.UUID   `bson:"form_id"`
	UnitID              uuid.UUID   `bson:"unit_id"`
	Name                string      `bson:"name"`
	ExpireTime          int         `bson:"expire_time"`
	EffectOnDriver      bool        `bson:"effect_on_driver"`
	EffectOnPregnant    bool        `bson:"effect_on_pregnant"`
	MethodOfApplication string      `bson:"method_of_application"`
	IsPrescription      bool        `bson:"is_prescription"`
	Substances          []Substance `bson:"substances"`
	Dosages             []Dosage    `bson:"dosages"`
	Contraindications   []uuid.UUID `bson:"contraindications"`
	Recommendations     []uuid.UUID `bson:"recommendations"`
}

type Dosage struct {
	ID                  uuid.UUID `bson:"id"`
	ValueFrom           int       `bson:"value_from"`
	ValueTo             int       `bson:"value_to"`
	DosageType          string    `bson:"dosage_type"`
	DosageValue         float32   `bson:"dosage_value"`
	NumberOfDosesPerDay int       `bson:"number_of_doses_per_day"`
}

type Substance struct {
	SubstanceID   uuid.UUID `bson:"substance_id"`
	Concentration float32   `bson:"concentration"`
}

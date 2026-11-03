package req

import "github.com/google/uuid"

type DosageRuleWriteDto struct {
	Type                string  `json:"type"`
	ValueFrom           int     `json:"value_from"`
	ValueTo             int     `json:"value_to"`
	DosageValue         float32 `json:"dosage_value"`
	NumberOfDosesPerDay int     `json:"number_of_doses_per_day"`
}

type MedicineWriteDto struct {
	Name                string               `json:"name" required:"true"`
	ExpireTime          int                  `json:"expire_time" required:"true"`
	IsPrescription      bool                 `json:"is_prescription" required:"true"`
	MethodOfApplication string               `json:"method_of_application" required:"true"`
	EffectOnPregnant    bool                 `json:"effect_on_pregnant" required:"true"`
	EffectOnDriver      bool                 `json:"effect_on_driver" required:"true"`
	FormID              uuid.UUID            `json:"form_id" required:"true"`
	UnitID              uuid.UUID            `json:"unit_id" required:"true"`
	Indications         []uuid.UUID          `json:"indications" required:"true"`
	Contraindications   []uuid.UUID          `json:"contraindications" required:"true"`
	Substances          []ActiveSubstanceDto `json:"substances" required:"true"`
	Dosages             []DosageRuleWriteDto `json:"dosages" required:"true"`
}

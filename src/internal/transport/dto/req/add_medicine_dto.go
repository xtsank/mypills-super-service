package req

import "github.com/google/uuid"

type AddMedicineDto struct {
	Name                string                `json:"name" binding:"required"`
	ExpireTime          int                   `json:"expire_time" binding:"required"`
	IsPrescription      bool                  `json:"is_prescription" binding:"required"`
	MethodOfApplication string                `json:"method_of_application" binding:"required"`
	EffectOnPregnant    bool                  `json:"effect_on_pregnant" binding:"required"`
	EffectOnDriver      bool                  `json:"effect_on_driver" binding:"required"`
	Form                uuid.UUID             `json:"form_id" binding:"required"`
	Unit                uuid.UUID             `json:"unit_id" binding:"required"`
	Recommendations     []uuid.UUID           `json:"recommendations" binding:"required"`
	Contraindications   []uuid.UUID           `json:"contraindications" binding:"required"`
	Substances          []*ActiveSubstanceDto `json:"substances" binding:"required"`
	Dosages             []*DosageRuleDto      `json:"dosages" binding:"required"`
}

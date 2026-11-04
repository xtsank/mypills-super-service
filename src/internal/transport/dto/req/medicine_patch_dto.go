package req

import "github.com/google/uuid"

type MedicinePatchDto struct {
	Name                *string               `json:"name"`
	ExpireTime          *int                  `json:"expire_time"`
	IsPrescription      *bool                 `json:"is_prescription"`
	MethodOfApplication *string               `json:"method_of_application"`
	EffectOnPregnant    *bool                 `json:"effect_on_pregnant"`
	EffectOnDriver      *bool                 `json:"effect_on_driver"`
	FormID              *uuid.UUID            `json:"form_id"`
	UnitID              *uuid.UUID            `json:"unit_id"`
	Indications         *[]uuid.UUID          `json:"indications"`
	Contraindications   *[]uuid.UUID          `json:"contraindications"`
	Substances          *[]ActiveSubstanceDto `json:"substances"`
	Dosages             *[]DosageRuleWriteDto `json:"dosages"`
}

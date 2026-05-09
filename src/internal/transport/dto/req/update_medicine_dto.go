package req

import "github.com/google/uuid"

type UpdateMedicineDto struct {
	ID                  uuid.UUID  `json:"id" binding:"required"`
	ExpireTime          *int       `json:"expire_time" binding:"required"`
	IsPrescription      *bool      `json:"is_prescription" binding:"required"`
	MethodOfApplication *string    `json:"method_of_application" binding:"required"`
	EffectOnPregnant    *bool      `json:"effect_on_pregnant" binding:"required"`
	EffectOnDriver      *bool      `json:"effect_on_driver" binding:"required"`
	FormID              *uuid.UUID `json:"form_id" binding:"required"`
	UnitID              *uuid.UUID `json:"unit_id" binding:"required"`
}

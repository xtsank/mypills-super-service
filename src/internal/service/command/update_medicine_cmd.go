package command

import "github.com/google/uuid"

type UpdateMedicineCmd struct {
	ID                  uuid.UUID
	ExpireTime          *int
	IsPrescription      *bool
	MethodOfApplication *string
	EffectOnPregnant    *bool
	EffectOnDriver      *bool
	FormID              *uuid.UUID
	UnitID              *uuid.UUID
}

func NewUpdateMedicineCmd(id uuid.UUID,
	expireTime *int,
	isPrescription *bool,
	methodOfApplication *string,
	effectOnPregnant *bool,
	effectOnDriver *bool,
	formId *uuid.UUID,
	unitId *uuid.UUID,
) *UpdateMedicineCmd {
	return &UpdateMedicineCmd{
		ID:                  id,
		ExpireTime:          expireTime,
		IsPrescription:      isPrescription,
		MethodOfApplication: methodOfApplication,
		EffectOnPregnant:    effectOnPregnant,
		EffectOnDriver:      effectOnDriver,
		FormID:              formId,
		UnitID:              unitId,
	}
}

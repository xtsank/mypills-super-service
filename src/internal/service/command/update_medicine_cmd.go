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

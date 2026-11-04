package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type MedicinePatchCmd struct {
	ID                  uuid.UUID
	Name                *string
	ExpireTime          *int
	IsPrescription      *bool
	MethodOfApplication *string
	EffectOnPregnant    *bool
	EffectOnDriver      *bool
	FormID              *uuid.UUID
	UnitID              *uuid.UUID
	Indications         *[]uuid.UUID
	Contraindications   *[]uuid.UUID
	Substances          *[]MedicineSubstance
	Dosages             *[]MedicineDosage
}

func NewMedicinePatchCmd(id uuid.UUID, input req.MedicinePatchDto) MedicinePatchCmd {
	cmd := MedicinePatchCmd{
		ID: id, Name: input.Name, ExpireTime: input.ExpireTime,
		IsPrescription: input.IsPrescription, MethodOfApplication: input.MethodOfApplication,
		EffectOnPregnant: input.EffectOnPregnant, EffectOnDriver: input.EffectOnDriver,
		FormID: input.FormID, UnitID: input.UnitID,
		Indications: input.Indications, Contraindications: input.Contraindications,
	}
	if input.Substances != nil {
		items := make([]MedicineSubstance, 0, len(*input.Substances))
		for _, item := range *input.Substances {
			items = append(items, MedicineSubstance{ID: item.ID, Concentration: item.Concentration})
		}
		cmd.Substances = &items
	}
	if input.Dosages != nil {
		items := make([]MedicineDosage, 0, len(*input.Dosages))
		for _, item := range *input.Dosages {
			items = append(items, MedicineDosage{
				Type: item.Type, ValueFrom: item.ValueFrom, ValueTo: item.ValueTo,
				DosageValue: item.DosageValue, NumberOfDosesPerDay: item.NumberOfDosesPerDay,
			})
		}
		cmd.Dosages = &items
	}
	return cmd
}

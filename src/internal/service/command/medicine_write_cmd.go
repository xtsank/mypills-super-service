package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type MedicineSubstance struct {
	ID            uuid.UUID
	Concentration float32
}

type MedicineDosage struct {
	Type                string
	ValueFrom           int
	ValueTo             int
	DosageValue         float32
	NumberOfDosesPerDay int
}

type MedicineWriteCmd struct {
	Name                string
	ExpireTime          int
	IsPrescription      bool
	MethodOfApplication string
	EffectOnPregnant    bool
	EffectOnDriver      bool
	FormID              uuid.UUID
	UnitID              uuid.UUID
	Indications         []uuid.UUID
	Contraindications   []uuid.UUID
	Substances          []MedicineSubstance
	Dosages             []MedicineDosage
}

func NewMedicineWriteCmd(input req.MedicineWriteDto) MedicineWriteCmd {
	substances := make([]MedicineSubstance, 0, len(input.Substances))
	for _, item := range input.Substances {
		substances = append(substances, MedicineSubstance{ID: item.ID, Concentration: item.Concentration})
	}
	dosages := make([]MedicineDosage, 0, len(input.Dosages))
	for _, item := range input.Dosages {
		dosages = append(dosages, MedicineDosage{
			Type: item.Type, ValueFrom: item.ValueFrom, ValueTo: item.ValueTo,
			DosageValue: item.DosageValue, NumberOfDosesPerDay: item.NumberOfDosesPerDay,
		})
	}
	return MedicineWriteCmd{
		Name: input.Name, ExpireTime: input.ExpireTime,
		IsPrescription: input.IsPrescription, MethodOfApplication: input.MethodOfApplication,
		EffectOnPregnant: input.EffectOnPregnant, EffectOnDriver: input.EffectOnDriver,
		FormID: input.FormID, UnitID: input.UnitID,
		Indications: input.Indications, Contraindications: input.Contraindications,
		Substances: substances, Dosages: dosages,
	}
}

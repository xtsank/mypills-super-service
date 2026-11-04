package service

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
)

func buildCatalogMedicine(id uuid.UUID, cmd command.MedicineWriteCmd) (*medicine.Medicine, error) {
	substances := make([]medicine.ActiveSubstance, 0, len(cmd.Substances))
	for _, item := range cmd.Substances {
		substances = append(substances, medicine.ActiveSubstance{ID: item.ID, Concentration: item.Concentration})
	}
	dosages := make([]medicine.DosageRule, 0, len(cmd.Dosages))
	for _, item := range cmd.Dosages {
		dosages = append(dosages, medicine.DosageRule{
			ID: uuid.New(), Type: medicine.DosageType(item.Type), ValueFrom: item.ValueFrom,
			ValueTo: item.ValueTo, DosageValue: item.DosageValue, NumberOfDosesPerDay: item.NumberOfDosesPerDay,
		})
	}
	return medicine.NewCatalogMedicine(id, cmd.Name, cmd.ExpireTime, cmd.IsPrescription,
		cmd.MethodOfApplication, cmd.EffectOnPregnant, cmd.EffectOnDriver,
		cmd.FormID, cmd.UnitID, substances, dosages, cmd.Contraindications, cmd.Indications)
}

func patchCatalogMedicine(current *medicine.Medicine, patch command.MedicinePatchCmd) (*medicine.Medicine, error) {
	cmd := command.MedicineWriteCmd{
		Name: current.Name, ExpireTime: current.ExpireTime,
		IsPrescription: current.IsPrescription, MethodOfApplication: current.MethodOfApplication,
		EffectOnPregnant: current.EffectOnPregnant, EffectOnDriver: current.EffectOnDriver,
		FormID: current.Form, UnitID: current.Unit,
		Indications:       append([]uuid.UUID{}, current.Recommendation...),
		Contraindications: append([]uuid.UUID{}, current.Contraindications...),
		Substances:        make([]command.MedicineSubstance, 0, len(current.Substances)),
		Dosages:           make([]command.MedicineDosage, 0, len(current.Dosages)),
	}
	for _, item := range current.Substances {
		cmd.Substances = append(cmd.Substances, command.MedicineSubstance{ID: item.ID, Concentration: item.Concentration})
	}
	for _, item := range current.Dosages {
		cmd.Dosages = append(cmd.Dosages, command.MedicineDosage{
			Type: string(item.Type), ValueFrom: item.ValueFrom, ValueTo: item.ValueTo,
			DosageValue: item.DosageValue, NumberOfDosesPerDay: item.NumberOfDosesPerDay,
		})
	}
	if patch.Name != nil {
		cmd.Name = *patch.Name
	}
	if patch.ExpireTime != nil {
		cmd.ExpireTime = *patch.ExpireTime
	}
	if patch.IsPrescription != nil {
		cmd.IsPrescription = *patch.IsPrescription
	}
	if patch.MethodOfApplication != nil {
		cmd.MethodOfApplication = *patch.MethodOfApplication
	}
	if patch.EffectOnPregnant != nil {
		cmd.EffectOnPregnant = *patch.EffectOnPregnant
	}
	if patch.EffectOnDriver != nil {
		cmd.EffectOnDriver = *patch.EffectOnDriver
	}
	if patch.FormID != nil {
		cmd.FormID = *patch.FormID
	}
	if patch.UnitID != nil {
		cmd.UnitID = *patch.UnitID
	}
	if patch.Indications != nil {
		cmd.Indications = *patch.Indications
	}
	if patch.Contraindications != nil {
		cmd.Contraindications = *patch.Contraindications
	}
	if patch.Substances != nil {
		cmd.Substances = *patch.Substances
	}
	if patch.Dosages != nil {
		cmd.Dosages = *patch.Dosages
	}
	updated, err := buildCatalogMedicine(current.ID, cmd)
	if err != nil {
		return nil, err
	}
	if patch.Dosages == nil {
		updated.Dosages = current.Dosages
	}
	return updated, nil
}

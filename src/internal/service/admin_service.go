package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IAdminService interface {
	AddMedicine(ctx context.Context, cmd *command.AddMedicineCmd) (*res.AdminResDto, error)
	UpdateMedicine(ctx context.Context, cmd *command.UpdateMedicineCmd) (*res.AdminResDto, error)
	RemoveMedicine(ctx context.Context, cmd *command.RemoveMedicineCmd) error
	UpdateIndications(ctx context.Context, cmd command.UpdateLinksCmd) error
	UpdateContraindications(ctx context.Context, cmd command.UpdateLinksCmd) error
	UpdateComposition(ctx context.Context, cmd command.UpdateCompositionCmd) error
	AddDosageRule(ctx context.Context, cmd command.AddDosageRuleCmd) error
	DeleteDosageRule(ctx context.Context, cmd command.RemoveDosageRuleCmd) error
}

type AdminService struct {
	medicineRepo medicine.IMedicineRepository
}

func NewAdminService(i do.Injector) (*AdminService, error) {
	medicineRepo := do.MustInvoke[medicine.IMedicineRepository](i)

	return &AdminService{medicineRepo: medicineRepo}, nil
}

func (s *AdminService) AddMedicine(ctx context.Context, cmd *command.AddMedicineCmd) (*res.AdminResDto, error) {
	var substances []medicine.ActiveSubstance
	for _, subDto := range cmd.Substances {
		substances = append(substances, medicine.ActiveSubstance{
			ID:            subDto.ID,
			Concentration: subDto.Concentration,
		})
	}

	var dosages []medicine.DosageRule
	for _, dDto := range cmd.Dosages {
		dosages = append(dosages, medicine.DosageRule{
			ID:                  uuid.New(),
			ValueFrom:           dDto.ValueFrom,
			ValueTo:             dDto.ValueTo,
			Type:                medicine.DosageType(dDto.Type),
			DosageValue:         dDto.DosageValue,
			NumberOfDosesPerDay: dDto.NumberOfDosesPerDay,
		})
	}

	newMed, err := medicine.NewMedicine(
		uuid.New(), cmd.Name, cmd.ExpireTime, cmd.IsPrescription,
		cmd.MethodOfApplication, cmd.EffectOnPregnant, cmd.EffectOnDriver,
		cmd.Form, cmd.Unit, substances, dosages,
		cmd.Contraindications, cmd.Recommendation,
	)
	if err != nil {
		return nil, err
	}

	err = s.medicineRepo.Create(ctx, newMed)
	if err != nil {
		return nil, err
	}

	return res.NewAdminResDto(newMed), nil
}

func (s *AdminService) UpdateMedicine(ctx context.Context, cmd *command.UpdateMedicineCmd) (*res.AdminResDto, error) {
	med, err := s.medicineRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if med == nil {
		return nil, errors.ErrMedicineNotFound
	}

	if cmd.ExpireTime != nil {
		med.ExpireTime = *cmd.ExpireTime
	}
	if cmd.IsPrescription != nil {
		med.IsPrescription = *cmd.IsPrescription
	}
	if cmd.MethodOfApplication != nil {
		med.MethodOfApplication = *cmd.MethodOfApplication
	}
	if cmd.EffectOnPregnant != nil {
		med.EffectOnPregnant = *cmd.EffectOnPregnant
	}
	if cmd.EffectOnDriver != nil {
		med.EffectOnDriver = *cmd.EffectOnDriver
	}
	if cmd.FormID != nil {
		med.Form = *cmd.FormID
	}
	if cmd.UnitID != nil {
		med.Unit = *cmd.UnitID
	}

	err = s.medicineRepo.Update(ctx, med)
	if err != nil {
		return nil, err
	}

	return res.NewAdminResDto(med), nil
}

func (s *AdminService) RemoveMedicine(ctx context.Context, cmd *command.RemoveMedicineCmd) error {
	return s.medicineRepo.Delete(ctx, cmd.ID)
}

func (s *AdminService) UpdateIndications(ctx context.Context, cmd command.UpdateLinksCmd) error {
	return s.medicineRepo.UpdateIndications(ctx, cmd.MedicineID, cmd.IDs)
}

func (s *AdminService) UpdateContraindications(ctx context.Context, cmd command.UpdateLinksCmd) error {
	return s.medicineRepo.UpdateContraindications(ctx, cmd.MedicineID, cmd.IDs)
}

func (s *AdminService) UpdateComposition(ctx context.Context, cmd command.UpdateCompositionCmd) error {
	var substances []medicine.ActiveSubstance
	for _, sub := range cmd.Substances {
		substances = append(substances, medicine.ActiveSubstance{
			ID:            sub.ID,
			Concentration: sub.Concentration,
		})
	}

	return s.medicineRepo.UpdateComposition(ctx, cmd.MedicineID, substances)
}

func (s *AdminService) AddDosageRule(ctx context.Context, cmd command.AddDosageRuleCmd) error {
	rule := &medicine.DosageRule{
		ID:                  uuid.New(),
		ValueFrom:           cmd.Dosage.ValueFrom,
		ValueTo:             cmd.Dosage.ValueTo,
		Type:                medicine.DosageType(cmd.Dosage.Type),
		DosageValue:         cmd.Dosage.DosageValue,
		NumberOfDosesPerDay: cmd.Dosage.NumberOfDosesPerDay,
	}

	return s.medicineRepo.AddDosageRule(ctx, cmd.MedicineID, rule)
}

func (s *AdminService) DeleteDosageRule(ctx context.Context, cmd command.RemoveDosageRuleCmd) error {
	return s.medicineRepo.DeleteDosageRule(ctx, cmd.RuleID)
}

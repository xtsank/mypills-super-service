package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type AddMedicineCmd struct {
	Name                string
	ExpireTime          int
	IsPrescription      bool
	MethodOfApplication string
	EffectOnPregnant    bool
	EffectOnDriver      bool
	Form                uuid.UUID
	Unit                uuid.UUID

	Recommendation    []uuid.UUID
	Contraindications []uuid.UUID

	Substances []*req.ActiveSubstanceDto
	Dosages    []*req.DosageRuleDto
}

func NewAddMedicineCmd(
	name string,
	expireTime int,
	isPrescription bool,
	method string,
	pregnantEffect bool,
	driverEffect bool,
	form uuid.UUID,
	unit uuid.UUID,
	recommendations []uuid.UUID,
	contraindications []uuid.UUID,
	substances []*req.ActiveSubstanceDto,
	dosages []*req.DosageRuleDto,
) (*AddMedicineCmd, error) {
	if name == "" {
		return nil, errors.ErrEmptyName.WithSource()
	}

	if expireTime <= 0 {
		return nil, errors.ErrExpireTimeTooLow.WithSource()
	}

	for _, s := range substances {
		if s.Concentration <= 0 {
			return nil, errors.ErrInvalidConcentration.WithSource()
		}
	}
	for _, d := range dosages {
		if d.ValueFrom < 0 || d.ValueTo < d.ValueFrom {
			return nil, errors.ErrInvalidDosageRange.WithSource()
		}
		if d.DosageValue <= 0 {
			return nil, errors.ErrInvalidDosageValue.WithSource()
		}
		if d.NumberOfDosesPerDay <= 0 {
			return nil, errors.ErrInvalidNumDoses.WithSource()
		}
	}

	return &AddMedicineCmd{
		Name:                name,
		ExpireTime:          expireTime,
		IsPrescription:      isPrescription,
		MethodOfApplication: method,
		EffectOnPregnant:    pregnantEffect,
		EffectOnDriver:      driverEffect,
		Form:                form,
		Unit:                unit,
		Recommendation:      recommendations,
		Contraindications:   contraindications,
		Substances:          substances,
		Dosages:             dosages,
	}, nil
}

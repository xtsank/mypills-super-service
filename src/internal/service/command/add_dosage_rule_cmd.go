package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type AddDosageRuleCmd struct {
	MedicineID uuid.UUID
	Dosage     *req.DosageRuleDto
}

func NewAddDosageRuleCmd(id uuid.UUID, dosage *req.DosageRuleDto) (*AddDosageRuleCmd, error) {
	if dosage.ValueFrom < 0 || dosage.ValueTo < dosage.ValueFrom {
		return nil, errors.ErrInvalidDosageRange.WithSource()
	}
	if dosage.DosageValue <= 0 {
		return nil, errors.ErrInvalidDosageValue.WithSource()
	}
	if dosage.NumberOfDosesPerDay <= 0 {
		return nil, errors.ErrInvalidNumDoses.WithSource()
	}

	return &AddDosageRuleCmd{
		MedicineID: id,
		Dosage:     dosage,
	}, nil
}

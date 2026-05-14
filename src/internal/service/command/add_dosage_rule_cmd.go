package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type AddDosageRuleCmd struct {
	MedicineID uuid.UUID
	Dosage     *req.DosageRuleDto
}

func NewAddDosageRuleCmd(id uuid.UUID, dosage *req.DosageRuleDto) (*AddDosageRuleCmd, error) {
	return &AddDosageRuleCmd{
		MedicineID: id,
		Dosage:     dosage,
	}, nil
}

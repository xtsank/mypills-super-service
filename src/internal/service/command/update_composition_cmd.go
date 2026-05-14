package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type UpdateCompositionCmd struct {
	MedicineID uuid.UUID
	Substances []*req.ActiveSubstanceDto
}

func NewUpdateCompositionCmd(id uuid.UUID, substances []*req.ActiveSubstanceDto) (*UpdateCompositionCmd, error) {
	for _, s := range substances {
		if s.Concentration <= 0 {
			return nil, errors.ErrInvalidConcentration.WithSource()
		}
	}

	return &UpdateCompositionCmd{
		MedicineID: id,
		Substances: substances,
	}, nil
}

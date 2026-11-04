package command

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type CreateCabinetItemCmd struct {
	UserID            uuid.UUID
	CabinetID         uuid.UUID
	MedicineID        uuid.UUID
	DateOfManufacture time.Time
	Quantity          float64
	MinimumQuantity   float64
}

func NewCreateCabinetItemCmd(userID, cabinetID uuid.UUID, input req.CabinetItemCreateDto) *CreateCabinetItemCmd {
	return &CreateCabinetItemCmd{UserID: userID, CabinetID: cabinetID, MedicineID: input.MedicineID,
		DateOfManufacture: input.DateOfManufacture.Time(), Quantity: input.Quantity, MinimumQuantity: input.MinimumQuantity}
}

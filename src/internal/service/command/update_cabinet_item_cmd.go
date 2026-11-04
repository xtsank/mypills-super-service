package command

import (
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/req"
)

type UpdateCabinetItemCmd struct {
	UserID            uuid.UUID
	CabinetID         uuid.UUID
	ItemID            uuid.UUID
	Quantity          *float64
	MinimumQuantity   *float64
	DateOfManufacture *time.Time
}

func NewUpdateCabinetItemCmd(userID, cabinetID, itemID uuid.UUID, input req.CabinetItemPatchDto) *UpdateCabinetItemCmd {
	var dateOfManufacture *time.Time
	if input.DateOfManufacture != nil {
		date := input.DateOfManufacture.Time()
		dateOfManufacture = &date
	}
	return &UpdateCabinetItemCmd{UserID: userID, CabinetID: cabinetID, ItemID: itemID,
		Quantity: input.Quantity, MinimumQuantity: input.MinimumQuantity, DateOfManufacture: dateOfManufacture}
}

package req

import "github.com/google/uuid"

type CabinetItemsListDto struct {
	ID                *uuid.UUID
	LowStock          bool
	ExpiresWithinDays *int
	Expired           bool
	Limit             int
	Offset            int
}

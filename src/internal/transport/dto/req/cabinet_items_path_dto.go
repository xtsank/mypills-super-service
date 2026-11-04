package req

import "github.com/google/uuid"

type CabinetItemsPathDto struct {
	UserID    uuid.UUID
	CabinetID uuid.UUID
}

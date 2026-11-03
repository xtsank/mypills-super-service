package req

import "github.com/google/uuid"

type MedicineListDto struct {
	ID     *uuid.UUID
	Name   *string
	Limit  int
	Offset int
}

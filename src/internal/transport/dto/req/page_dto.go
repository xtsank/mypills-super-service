package req

import "github.com/google/uuid"

type PageDto struct {
	ID     *uuid.UUID
	Limit  int
	Offset int
}

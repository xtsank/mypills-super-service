package cabinet

import (
	"context"

	"github.com/google/uuid"
)

type ICabinetRepository interface {
	EnsureLegacy(ctx context.Context, userID uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID, id *uuid.UUID, limit, offset int) (Page, error)
	Find(ctx context.Context, userID, id uuid.UUID) (*Cabinet, error)
	Create(ctx context.Context, item Cabinet) error
	Rename(ctx context.Context, item Cabinet) error
	Delete(ctx context.Context, userID, id uuid.UUID) error
}

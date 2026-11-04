package cabinet

import (
	"context"

	"github.com/google/uuid"
)

type ICabinetItemRepository interface {
	List(ctx context.Context, userID, cabinetID uuid.UUID) ([]Item, error)
	Find(ctx context.Context, userID, cabinetID, id uuid.UUID) (*Item, error)
	Create(ctx context.Context, item Item) error
	Update(ctx context.Context, item Item) error
	Delete(ctx context.Context, userID, cabinetID, id uuid.UUID) error
}

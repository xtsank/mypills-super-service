package restock

import (
	"context"

	"github.com/google/uuid"
)

type IRestockRepository interface {
	List(ctx context.Context, userID uuid.UUID, id *uuid.UUID, limit, offset int) (Page, error)
	Find(ctx context.Context, userID, id uuid.UUID) (*Item, error)
	Create(ctx context.Context, item Item) error
	Update(ctx context.Context, item Item) error
	Delete(ctx context.Context, userID, id uuid.UUID) error
}

package medicine

import (
	"context"

	"github.com/google/uuid"
)

type ICatalogRepository interface {
	List(ctx context.Context, id *uuid.UUID, name string, limit, offset int) (Page, error)
}

package dictionary

import "context"

type IRepository interface {
	ListIllnesses(ctx context.Context) ([]Item, error)
	ListSubstances(ctx context.Context) ([]Item, error)
	ListForms(ctx context.Context) ([]Item, error)
	ListUnits(ctx context.Context) ([]Item, error)
	ListMedicines(ctx context.Context) ([]Item, error)
}

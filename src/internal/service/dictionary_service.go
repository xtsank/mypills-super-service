package service

import (
	"context"
	"sort"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/dictionary"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IDictionaryService interface {
	ReferenceIndex(ctx context.Context, cmd *command.GetReferenceIndexCmd) (*res.ReferenceIndexDto, error)
	ReferenceList(ctx context.Context, cmd *command.ListReferenceCmd) (*res.ReferencePageDto, error)
	ListIllnesses(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListSubstances(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListForms(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListUnits(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListMedicines(ctx context.Context) ([]res.DictionaryItemDto, error)
}

func (s *DictionaryService) ReferenceIndex(ctx context.Context, cmd *command.GetReferenceIndexCmd) (*res.ReferenceIndexDto, error) {
	return res.NewReferenceIndexDto(), nil
}

func (s *DictionaryService) ReferenceList(ctx context.Context, cmd *command.ListReferenceCmd) (*res.ReferencePageDto, error) {
	var items []res.DictionaryItemDto
	var err error
	switch cmd.Kind {
	case "illnesses":
		items, err = s.ListIllnesses(ctx)
	case "substances":
		items, err = s.ListSubstances(ctx)
	case "forms":
		items, err = s.ListForms(ctx)
	case "units":
		items, err = s.ListUnits(ctx)
	default:
		return nil, apperrors.ErrReferenceNotFound.WithSource()
	}
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Name == items[j].Name {
			return items[i].ID.String() < items[j].ID.String()
		}
		return items[i].Name < items[j].Name
	})
	return res.NewReferencePageDto(items, cmd.Limit, cmd.Offset), nil
}

type DictionaryService struct {
	repo dictionary.IRepository
}

func NewDictionaryService(i do.Injector) (IDictionaryService, error) {
	repo := do.MustInvoke[dictionary.IRepository](i)
	return &DictionaryService{repo: repo}, nil
}

func (s *DictionaryService) ListIllnesses(ctx context.Context) ([]res.DictionaryItemDto, error) {
	items, err := s.repo.ListIllnesses(ctx)
	return res.NewDictionaryItems(items), err
}

func (s *DictionaryService) ListSubstances(ctx context.Context) ([]res.DictionaryItemDto, error) {
	items, err := s.repo.ListSubstances(ctx)
	return res.NewDictionaryItems(items), err
}

func (s *DictionaryService) ListForms(ctx context.Context) ([]res.DictionaryItemDto, error) {
	items, err := s.repo.ListForms(ctx)
	return res.NewDictionaryItems(items), err
}

func (s *DictionaryService) ListUnits(ctx context.Context) ([]res.DictionaryItemDto, error) {
	items, err := s.repo.ListUnits(ctx)
	return res.NewDictionaryItems(items), err
}

func (s *DictionaryService) ListMedicines(ctx context.Context) ([]res.DictionaryItemDto, error) {
	items, err := s.repo.ListMedicines(ctx)
	return res.NewDictionaryItems(items), err
}

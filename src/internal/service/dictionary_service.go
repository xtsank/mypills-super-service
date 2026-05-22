package service

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/dictionary"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IDictionaryService interface {
	ListIllnesses(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListSubstances(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListForms(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListUnits(ctx context.Context) ([]res.DictionaryItemDto, error)
	ListMedicines(ctx context.Context) ([]res.DictionaryItemDto, error)
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


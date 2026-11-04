package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IMedicineCatalogService interface {
	List(ctx context.Context, cmd *command.ListMedicinesCmd) (*res.MedicineCatalogPageDto, error)
	Create(ctx context.Context, cmd command.MedicineWriteCmd) (*res.MedicineCatalogDto, error)
	Update(ctx context.Context, cmd command.MedicinePatchCmd) (*res.MedicineCatalogDto, error)
	Delete(ctx context.Context, cmd *command.DeleteMedicineCmd) error
}

type MedicineCatalogService struct {
	repo    medicine.IMedicineRepository
	catalog medicine.ICatalogRepository
}

func NewMedicineCatalogService(i do.Injector) (IMedicineCatalogService, error) {
	return &MedicineCatalogService{
		repo:    do.MustInvoke[medicine.IMedicineRepository](i),
		catalog: do.MustInvoke[medicine.ICatalogRepository](i),
	}, nil
}

func (s *MedicineCatalogService) List(ctx context.Context, cmd *command.ListMedicinesCmd) (*res.MedicineCatalogPageDto, error) {
	name := ""
	if cmd.Name != nil {
		name = strings.TrimSpace(*cmd.Name)
		if name == "" {
			return nil, apperrors.ErrInvalidInput.WithSource()
		}
	}
	if len([]rune(name)) > 100 {
		return nil, apperrors.ErrInvalidInput.WithSource()
	}
	page, err := s.catalog.List(ctx, cmd.ID, name, cmd.Limit, cmd.Offset)
	if err != nil {
		return nil, err
	}
	return res.NewMedicineCatalogPageDto(page), nil
}

func (s *MedicineCatalogService) Create(ctx context.Context, cmd command.MedicineWriteCmd) (*res.MedicineCatalogDto, error) {
	item, err := buildCatalogMedicine(uuid.New(), cmd)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return res.NewMedicineCatalogDto(item), nil
}

func (s *MedicineCatalogService) Update(ctx context.Context, cmd command.MedicinePatchCmd) (*res.MedicineCatalogDto, error) {
	current, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, apperrors.ErrMedicineNotFound.WithSource()
	}
	updated, err := patchCatalogMedicine(current, cmd)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, updated); err != nil {
		return nil, err
	}
	return res.NewMedicineCatalogDto(updated), nil
}

func (s *MedicineCatalogService) Delete(ctx context.Context, cmd *command.DeleteMedicineCmd) error {
	if _, err := s.repo.FindByID(ctx, cmd.ID); err != nil {
		return err
	}
	return s.repo.Delete(ctx, cmd.ID)
}

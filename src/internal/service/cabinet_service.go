package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type ICabinetService interface {
	AddItem(ctx context.Context, cmd *command.AddItemCmd) (*res.CabinetResDto, error)
	RemoveItem(ctx context.Context, cmd *command.RemoveItemCmd) error
	UpdateQty(ctx context.Context, cmd *command.UpdateQtyCmd) (*res.CabinetResDto, error)
}

type CabinetService struct {
	cabinetRepo cabinet_item.ICabinetItemRepository
}

func NewCabinetService(i do.Injector) (*CabinetService, error) {
	cabinetRepo := do.MustInvoke[cabinet_item.ICabinetItemRepository](i)

	return &CabinetService{cabinetRepo: cabinetRepo}, nil
}

func (s *CabinetService) AddItem(ctx context.Context, cmd *command.AddItemCmd) (*res.CabinetResDto, error) {
	existingItem, err := s.cabinetRepo.FindExistingCabinetItem(ctx, cmd.UserID, cmd.MedicineID, cmd.DateOfManufacture)
	if err != nil {
		return nil, err
	}
	if existingItem != nil {
		existingItem.Quantity += cmd.Quantity
		err = s.cabinetRepo.Update(ctx, existingItem)
		if err != nil {
			return nil, err
		}
		return res.NewCabinetResDto(existingItem), nil
	}

	item, err := cabinet_item.NewCabinetItem(
		uuid.New(),
		cmd.UserID,
		cmd.MedicineID,
		cmd.DateOfManufacture,
		cmd.Quantity,
	)
	if err != nil {
		return nil, err
	}

	err = s.cabinetRepo.Save(ctx, item)
	if err != nil {
		return nil, err
	}

	return res.NewCabinetResDto(item), nil
}

func (s *CabinetService) RemoveItem(ctx context.Context, cmd *command.RemoveItemCmd) error {
	return s.cabinetRepo.Delete(ctx, cmd.ID)
}

func (s *CabinetService) UpdateQty(ctx context.Context, cmd *command.UpdateQtyCmd) (*res.CabinetResDto, error) {
	existingItem, err := s.cabinetRepo.FindById(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if existingItem == nil {
		return nil, errors.ErrCabinetItemNotFound
	}

	existingItem.Quantity = cmd.Quantity
	err = s.cabinetRepo.Update(ctx, existingItem)
	if err != nil {
		return nil, err
	}

	return res.NewCabinetResDto(existingItem), nil
}

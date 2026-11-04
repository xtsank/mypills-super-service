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

type ILegacyCabinetService interface {
	AddItem(ctx context.Context, cmd *command.AddItemCmd) (*res.CabinetResDto, error)
	RemoveItem(ctx context.Context, cmd *command.RemoveItemCmd) error
	UpdateQty(ctx context.Context, cmd *command.UpdateQtyCmd) (*res.CabinetResDto, error)
	ListItems(ctx context.Context, userID uuid.UUID) ([]res.CabinetItemDetailsResDto, error)
}

type LegacyCabinetService struct {
	cabinetRepo cabinet_item.ICabinetItemRepository
}

func NewLegacyCabinetService(i do.Injector) (ILegacyCabinetService, error) {
	cabinetRepo := do.MustInvoke[cabinet_item.ICabinetItemRepository](i)

	return &LegacyCabinetService{cabinetRepo: cabinetRepo}, nil
}

func (s *LegacyCabinetService) AddItem(ctx context.Context, cmd *command.AddItemCmd) (*res.CabinetResDto, error) {
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

func (s *LegacyCabinetService) RemoveItem(ctx context.Context, cmd *command.RemoveItemCmd) error {
	if cmd.UserID != uuid.Nil {
		item, err := s.cabinetRepo.FindById(ctx, cmd.ID)
		if err != nil {
			return err
		}
		if item == nil || item.UserID != cmd.UserID {
			return errors.New("CABINET_ITEM_NOT_FOUND", "item not found", 404)
		}
	}
	return s.cabinetRepo.Delete(ctx, cmd.ID)
}

func (s *LegacyCabinetService) UpdateQty(ctx context.Context, cmd *command.UpdateQtyCmd) (*res.CabinetResDto, error) {
	existingItem, err := s.cabinetRepo.FindById(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if existingItem == nil || (cmd.UserID != uuid.Nil && existingItem.UserID != cmd.UserID) {
		return nil, errors.New("CABINET_ITEM_NOT_FOUND", "item not found", 404)
	}

	existingItem.Quantity = cmd.Quantity
	err = s.cabinetRepo.Update(ctx, existingItem)
	if err != nil {
		return nil, err
	}

	return res.NewCabinetResDto(existingItem), nil
}

func (s *LegacyCabinetService) ListItems(ctx context.Context, userID uuid.UUID) ([]res.CabinetItemDetailsResDto, error) {
	items, err := s.cabinetRepo.FindDetailsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return res.NewCabinetItemDetailsResDto(items), nil
}

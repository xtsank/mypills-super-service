package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"github.com/xtsank/mypills-super-service/src/internal/domain/restock"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IRestockService interface {
	List(ctx context.Context, cmd *command.ListRestockCmd) (*res.RestockPageDto, error)
	Create(ctx context.Context, cmd *command.CreateRestockCmd) (*res.RestockItemDto, error)
	Update(ctx context.Context, cmd *command.UpdateRestockCmd) (*res.RestockItemDto, error)
	Delete(ctx context.Context, cmd *command.DeleteRestockCmd) error
}

type RestockService struct {
	repo      restock.IRestockRepository
	medicines medicine.IMedicineRepository
}

func NewRestockService(i do.Injector) (IRestockService, error) {
	return &RestockService{repo: do.MustInvoke[restock.IRestockRepository](i),
		medicines: do.MustInvoke[medicine.IMedicineRepository](i)}, nil
}

func (s *RestockService) List(ctx context.Context, cmd *command.ListRestockCmd) (*res.RestockPageDto, error) {
	page, err := s.repo.List(ctx, cmd.UserID, cmd.ID, cmd.Limit, cmd.Offset)
	if err != nil {
		return nil, err
	}
	return res.NewRestockPageDto(page), nil
}

func (s *RestockService) Create(ctx context.Context, cmd *command.CreateRestockCmd) (*res.RestockItemDto, error) {
	med, err := s.medicines.FindByID(ctx, cmd.MedicineID)
	if err != nil {
		return nil, err
	}
	if med == nil {
		return nil, apperrors.ErrMedicineNotFound.WithSource()
	}
	item, err := restock.NewItem(uuid.New(), cmd.UserID, cmd.MedicineID, med.Name, cmd.DesiredQuantity)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, *item); err != nil {
		return nil, err
	}
	return res.NewRestockItemDto(*item), nil
}

func (s *RestockService) Update(ctx context.Context, cmd *command.UpdateRestockCmd) (*res.RestockItemDto, error) {
	item, err := s.repo.Find(ctx, cmd.UserID, cmd.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, apperrors.ErrRestockItemNotFound.WithSource()
	}
	updated, err := restock.NewItem(item.ID, item.UserID, item.MedicineID, item.MedicineName, cmd.DesiredQuantity)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Update(ctx, *updated); err != nil {
		return nil, err
	}
	return res.NewRestockItemDto(*updated), nil
}

func (s *RestockService) Delete(ctx context.Context, cmd *command.DeleteRestockCmd) error {
	item, err := s.repo.Find(ctx, cmd.UserID, cmd.ItemID)
	if err != nil {
		return err
	}
	if item == nil {
		return apperrors.ErrRestockItemNotFound.WithSource()
	}
	return s.repo.Delete(ctx, cmd.UserID, cmd.ItemID)
}

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type ICabinetService interface {
	List(ctx context.Context, cmd *command.ListCabinetsCmd) (*res.CabinetsPageDto, error)
	Create(ctx context.Context, cmd *command.CreateCabinetCmd) (*res.CabinetDto, error)
	Rename(ctx context.Context, cmd *command.RenameCabinetCmd) (*res.CabinetDto, error)
	Delete(ctx context.Context, cmd *command.DeleteCabinetCmd) error
}

type CabinetService struct{ repo cabinet.ICabinetRepository }

func NewCabinetService(i do.Injector) (ICabinetService, error) {
	return &CabinetService{repo: do.MustInvoke[cabinet.ICabinetRepository](i)}, nil
}

func (s *CabinetService) List(ctx context.Context, cmd *command.ListCabinetsCmd) (*res.CabinetsPageDto, error) {
	if err := s.repo.EnsureLegacy(ctx, cmd.UserID); err != nil {
		return nil, err
	}
	page, err := s.repo.List(ctx, cmd.UserID, cmd.ID, cmd.Limit, cmd.Offset)
	if err != nil {
		return nil, err
	}
	return res.NewCabinetsPageDto(page), nil
}

func (s *CabinetService) Create(ctx context.Context, cmd *command.CreateCabinetCmd) (*res.CabinetDto, error) {
	item, err := cabinet.NewCabinet(uuid.New(), cmd.UserID, cmd.Name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, *item); err != nil {
		return nil, err
	}
	return res.NewCabinetDto(*item), nil
}

func (s *CabinetService) Rename(ctx context.Context, cmd *command.RenameCabinetCmd) (*res.CabinetDto, error) {
	item, err := s.repo.Find(ctx, cmd.UserID, cmd.CabinetID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, apperrors.ErrCabinetNotFound.WithSource()
	}
	updated, err := cabinet.NewCabinet(item.ID, item.UserID, cmd.Name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Rename(ctx, *updated); err != nil {
		return nil, err
	}
	return res.NewCabinetDto(*updated), nil
}

func (s *CabinetService) Delete(ctx context.Context, cmd *command.DeleteCabinetCmd) error {
	item, err := s.repo.Find(ctx, cmd.UserID, cmd.CabinetID)
	if err != nil {
		return err
	}
	if item == nil {
		return apperrors.ErrCabinetNotFound.WithSource()
	}
	return s.repo.Delete(ctx, cmd.UserID, cmd.CabinetID)
}

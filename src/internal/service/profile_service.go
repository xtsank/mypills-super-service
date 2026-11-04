package service

import (
	"context"
	stderrors "errors"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IProfileService interface {
	GetProfile(ctx context.Context, cmd *command.GetProfileCmd) (*res.ProfileResDto, error)
	UpdateProfile(ctx context.Context, cmd *command.UpdateProfileCmd) (*res.ProfileResDto, error)
}

type ProfileService struct {
	userRepo user.IUserRepository
}

func NewProfileService(i do.Injector) (IProfileService, error) {
	userRepo := do.MustInvoke[user.IUserRepository](i)

	return &ProfileService{userRepo: userRepo}, nil
}

func (s *ProfileService) GetProfile(ctx context.Context, cmd *command.GetProfileCmd) (*res.ProfileResDto, error) {
	u, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if stderrors.Is(err, errors.ErrUserNotFound) {
		return nil, errors.ErrUnauthorized.WithSource()
	}
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.ErrUnauthorized.WithSource()
	}
	return res.NewProfileResDto(u), nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, cmd *command.UpdateProfileCmd) (*res.ProfileResDto, error) {
	u, err := s.userRepo.FindByID(ctx, cmd.ID)
	if stderrors.Is(err, errors.ErrUserNotFound) {
		return nil, errors.ErrUnauthorized.WithSource()
	}
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.ErrUnauthorized.WithSource()
	}

	updated, err := user.NewUpdatedUser(u, cmd.Email, cmd.Sex, cmd.Weight, cmd.Age,
		cmd.IsPregnant, cmd.IsDriver, cmd.Illnesses, cmd.Allergies)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Update(ctx, updated)
	if err != nil {
		return nil, err
	}

	return res.NewProfileResDto(updated), nil
}

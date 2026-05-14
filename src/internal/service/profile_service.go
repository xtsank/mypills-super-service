package service

import (
	"context"

	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IProfileService interface {
	UpdateProfile(ctx context.Context, cmd *command.UpdateProfileCmd) (*res.ProfileResDto, error)
}

type ProfileService struct {
	userRepo user.IUserRepository
}

func NewProfileService(i do.Injector) (IProfileService, error) {
	userRepo := do.MustInvoke[user.IUserRepository](i)

	return &ProfileService{userRepo: userRepo}, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, cmd *command.UpdateProfileCmd) (*res.ProfileResDto, error) {
	u, err := s.userRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.ErrUserNotFound.WithSource()
	}

	if cmd.Sex != nil {
		u.Sex = *cmd.Sex
	}
	if cmd.Weight != nil {
		u.Weight = *cmd.Weight
	}
	if cmd.Age != nil {
		u.Age = *cmd.Age
	}
	if cmd.IsPregnant != nil {
		u.IsPregnant = *cmd.IsPregnant
	}
	if cmd.IsDriver != nil {
		u.IsDriver = *cmd.IsDriver
	}
	if cmd.Illnesses != nil {
		u.Illnesses = cmd.Illnesses
	}
	if cmd.Allergies != nil {
		u.Allergies = cmd.Allergies
	}

	err = s.userRepo.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return res.NewProfileResDto(u), nil
}

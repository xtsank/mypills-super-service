package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
	"github.com/xtsank/mypills-super-service/src/internal/transport/dto/res"
)

type IAuthService interface {
	Register(ctx context.Context, cmd *command.CreateUserCmd) (*res.AuthResDto, error)
	Login(ctx context.Context, cmd *command.LoginUserCmd) (*res.AuthResDto, error)
}

type AuthService struct {
	userRepo     user.IUserRepository
	hasher       IPasswordHasher
	tokenManager TokenManager
}

func NewAuthService(i do.Injector) (IAuthService, error) {
	userRepo := do.MustInvoke[user.IUserRepository](i)
	hasher := do.MustInvoke[IPasswordHasher](i)
	tokenManager := do.MustInvoke[TokenManager](i)

	return &AuthService{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenManager: tokenManager,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, cmd *command.CreateUserCmd) (*res.AuthResDto, error) {
	exists, err := s.userRepo.ExistsByLogin(ctx, cmd.Login)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.ErrUserExists
	}

	hashedPassword, err := s.hasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	u, err := user.NewUser(
		uuid.New(),
		cmd.Login,
		hashedPassword,
		cmd.IsAdmin,
		cmd.Sex,
		cmd.Weight,
		cmd.Age,
		cmd.IsPregnant,
		cmd.IsDriver,
		cmd.Illnesses,
		cmd.Allergies,
	)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	token, err := s.tokenManager.GenerateToken(u.ID, u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return res.NewAuthResDto(u, token), nil
}

func (s *AuthService) Login(ctx context.Context, cmd *command.LoginUserCmd) (*res.AuthResDto, error) {
	u, err := s.userRepo.FindByLogin(ctx, cmd.Login)
	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.ErrUserNotFound
	}

	err = s.hasher.Compare(u.Password, cmd.Password)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	token, err := s.tokenManager.GenerateToken(u.ID, u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return res.NewAuthResDto(u, token), nil
}

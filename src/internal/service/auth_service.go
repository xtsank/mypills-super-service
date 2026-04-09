package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/internal/domain/user"
	"github.com/xtsank/mypills-super-service/internal/service/command"
)

type IAuthService interface {
	Register(ctx context.Context, cmd command.CreateUserCmd) (*user.User, error)
}

type AuthService struct {
	userRepo user.IUserRepository
}

func NewAuthService(userRepo user.IUserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, cmd command.CreateUserCmd) (*user.User, error) {
	existingUser, err := s.userRepo.FindByLogin(ctx, cmd.Login)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword := cmd.Password //TODO

	newUser := user.User{
		ID:         uuid.New(),
		Login:      cmd.Login,
		Password:   hashedPassword,
		Sex:        cmd.Sex,
		Weight:     cmd.Weight,
		Age:        cmd.Age,
		IsPregnant: cmd.IsPregnant,
		IsDriver:   cmd.IsDriver,
		Illnesses:  cmd.Illnesses,
		Allergies:  cmd.Allergies,
	}

	err = s.userRepo.Save(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

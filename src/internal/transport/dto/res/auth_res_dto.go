package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

type AuthResDto struct {
	Token string         `json:"token"`
	User  *ProfileResDto `json:"user"`
}

func NewAuthResDto(u *user.User, token string) *AuthResDto {
	return &AuthResDto{
		Token: uuid.New().String(),
		User:  NewProfileResDto(u),
	}
}

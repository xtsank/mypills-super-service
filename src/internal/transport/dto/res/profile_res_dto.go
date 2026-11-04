package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

type ProfileResDto struct {
	ID         string      `json:"id"`
	Login      string      `json:"login"`
	Email      string      `json:"email"`
	Sex        bool        `json:"sex"`
	Weight     int         `json:"weight"`
	Age        int         `json:"age"`
	IsPregnant bool        `json:"is_pregnant"`
	IsDriver   bool        `json:"is_driver"`
	Illnesses  []uuid.UUID `json:"illnesses"`
	Allergies  []uuid.UUID `json:"allergies"`
}

func NewProfileResDto(u *user.User) *ProfileResDto {
	illnesses := u.Illnesses
	if illnesses == nil {
		illnesses = []uuid.UUID{}
	}
	allergies := u.Allergies
	if allergies == nil {
		allergies = []uuid.UUID{}
	}
	return &ProfileResDto{
		ID:         u.ID.String(),
		Login:      u.Login,
		Email:      u.Email,
		Sex:        u.Sex,
		Weight:     u.Weight,
		Age:        u.Age,
		IsPregnant: u.IsPregnant,
		IsDriver:   u.IsDriver,
		Illnesses:  illnesses,
		Allergies:  allergies,
	}
}

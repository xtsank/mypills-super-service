package res

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

type ProfileResDto struct {
	ID         string      `json:"id"`
	Login      string      `json:"login"`
	Sex        bool        `json:"sex"`
	Weight     int         `json:"weight"`
	Age        int         `json:"age"`
	IsPregnant bool        `json:"is_pregnant"`
	IsDriver   bool        `json:"is_driver"`
	Illnesses  []uuid.UUID `json:"illnesses"`
	Allergies  []uuid.UUID `json:"allergies"`
}

func NewProfileResDto(u *user.User) *ProfileResDto {
	return &ProfileResDto{
		ID:         u.ID.String(),
		Login:      u.Login,
		Sex:        u.Sex,
		Age:        u.Age,
		IsPregnant: u.IsPregnant,
		IsDriver:   u.IsDriver,
		Illnesses:  u.Illnesses,
		Allergies:  u.Allergies,
	}
}

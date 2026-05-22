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
	NotifyEnabled         bool `json:"notify_enabled"`
	NotifyIntervalMinutes int  `json:"notify_interval_minutes"`
	Illnesses  []uuid.UUID `json:"illnesses"`
	Allergies  []uuid.UUID `json:"allergies"`
}

func NewProfileResDto(u *user.User) *ProfileResDto {
	return &ProfileResDto{
		ID:         u.ID.String(),
		Login:      u.Login,
		Email:      u.Email,
		Sex:        u.Sex,
		Weight:     u.Weight,
		Age:        u.Age,
		IsPregnant: u.IsPregnant,
		IsDriver:   u.IsDriver,
		NotifyEnabled:         u.Notify.Enabled,
		NotifyIntervalMinutes: u.Notify.IntervalMinutes,
		Illnesses:  u.Illnesses,
		Allergies:  u.Allergies,
	}
}

package req

import "github.com/google/uuid"

type CreateUserDto struct {
	Login      string      `json:"login" binding:"required,min=8"`
	Password   string      `json:"password" binding:"required,min=8"`
	Sex        bool        `json:"sex"`
	Weight     int         `json:"weight" binding:"required"`
	Age        int         `json:"age" binding:"required"`
	IsPregnant bool        `json:"is_pregnant"`
	IsDriver   bool        `json:"is_driver"`
	Illnesses  []uuid.UUID `json:"illnesses,omitempty"`
	Allergies  []uuid.UUID `json:"allergies,omitempty"`
}

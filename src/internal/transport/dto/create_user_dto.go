package dto

import "github.com/google/uuid"

type CreateUserDto struct {
	Login      string      `json:"login" binding:"required"`
	Password   string      `json:"password" binding:"required"`
	Sex        bool        `json:"sex"`
	Weight     float32     `json:"weight" binding:"required"`
	Age        int         `json:"age" binding:"required"`
	IsPregnant bool        `json:"is_pregnant"`
	IsDriver   bool        `json:"is_driver"`
	Illnesses  []uuid.UUID `json:"illnesses,omitempty"`
	Allergies  []uuid.UUID `json:"allergies,omitempty"`
}

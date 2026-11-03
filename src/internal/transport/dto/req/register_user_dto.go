package req

import "github.com/google/uuid"

type RegisterUserDto struct {
	Login      string      `json:"login" required:"true"`
	Email      string      `json:"email" required:"true"`
	Password   string      `json:"password" required:"true"`
	Age        int         `json:"age" required:"true"`
	Weight     int         `json:"weight" required:"true"`
	Sex        bool        `json:"sex"`
	IsPregnant bool        `json:"is_pregnant"`
	IsDriver   bool        `json:"is_driver"`
	Allergies  []uuid.UUID `json:"allergies"`
	Illnesses  []uuid.UUID `json:"illnesses"`
}

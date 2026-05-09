package command

import (
	"github.com/google/uuid"
)

type CreateUserCmd struct {
	Login      string
	Password   string
	IsAdmin    bool
	Sex        bool
	Weight     int
	Age        int
	IsPregnant bool
	IsDriver   bool
	Illnesses  []uuid.UUID
	Allergies  []uuid.UUID
}

func NewCreateUserCmd(
	login string,
	password string,
	sex bool,
	weight int,
	age int,
	isPregnant bool,
	isDriver bool,
	illnesses []uuid.UUID,
	allergies []uuid.UUID,
) *CreateUserCmd {
	return &CreateUserCmd{
		Login:      login,
		Password:   password,
		IsAdmin:    false,
		Sex:        sex,
		Weight:     weight,
		Age:        age,
		IsPregnant: isPregnant,
		IsDriver:   isDriver,
		Illnesses:  illnesses,
		Allergies:  allergies,
	}
}

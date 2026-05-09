package command

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
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
) (*CreateUserCmd, error) {
	if len(login) < 8 {
		return nil, errors.ErrLoginTooShort
	}

	if len(password) < 8 {
		return nil, errors.ErrPasswordTooShort
	}

	if weight <= 0 || weight > 500 {
		return nil, errors.ErrWrongWeight
	}

	if age <= 0 || age > 120 {
		return nil, errors.ErrWrongAge
	}

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
	}, nil
}

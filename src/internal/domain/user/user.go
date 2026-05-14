package user

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

type User struct {
	ID         uuid.UUID
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

func NewUser(
	id uuid.UUID,
	login string,
	password string,
	isAdmin bool,
	sex bool,
	weight int,
	age int,
	isPregnant bool,
	isDriver bool,
	illnesses []uuid.UUID,
	allergies []uuid.UUID,
) (*User, error) {
	if len(login) == 0 {
		return nil, errors.ErrLoginTooShort.WithSource()
	}

	if len(password) == 0 {
		return nil, errors.ErrPasswordTooShort.WithSource()
	}

	if weight <= 0 || weight > 500 {
		return nil, errors.ErrWrongWeight.WithSource()
	}

	if age <= 0 || age > 120 {
		return nil, errors.ErrWrongAge.WithSource()
	}

	if illnesses == nil {
		illnesses = []uuid.UUID{}
	}
	if allergies == nil {
		allergies = []uuid.UUID{}
	}

	return &User{
		ID:         id,
		Login:      login,
		Password:   password,
		IsAdmin:    isAdmin,
		Sex:        sex,
		Weight:     weight,
		Age:        age,
		IsPregnant: isPregnant,
		IsDriver:   isDriver,
		Illnesses:  illnesses,
		Allergies:  allergies,
	}, nil
}

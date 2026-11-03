package user

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

type User struct {
	ID         uuid.UUID
	Login      string
	Email      string
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
	email string,
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

	if len(email) == 0 {
		return nil, errors.ErrEmailRequired.WithSource()
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
		Email:      email,
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

func (u *User) UpdateProfile(email *string, sex *bool, weight, age *int,
	isPregnant, isDriver *bool, illnesses, allergies []uuid.UUID) error {
	if email != nil && !validEmail(*email) {
		return errors.ErrEmailRequired.WithSource()
	}
	if weight != nil && (*weight < 1 || *weight > 500) {
		return errors.ErrWrongWeight.WithSource()
	}
	if age != nil && (*age < 1 || *age > 120) {
		return errors.ErrWrongAge.WithSource()
	}
	if !uniqueIDs(illnesses) || !uniqueIDs(allergies) {
		return errors.ErrInvalidInput.WithSource()
	}
	if email != nil {
		u.Email = *email
	}
	if weight != nil {
		u.Weight = *weight
	}
	if age != nil {
		u.Age = *age
	}
	if sex != nil {
		u.Sex = *sex
	}
	if isPregnant != nil {
		u.IsPregnant = *isPregnant
	}
	if isDriver != nil {
		u.IsDriver = *isDriver
	}
	if illnesses != nil {
		u.Illnesses = illnesses
	}
	if allergies != nil {
		u.Allergies = allergies
	}
	return nil
}

func NewUpdatedUser(current *User, email *string, sex *bool, weight, age *int,
	isPregnant, isDriver *bool, illnesses, allergies []uuid.UUID) (*User, error) {
	updated := *current
	if err := updated.UpdateProfile(email, sex, weight, age, isPregnant, isDriver, illnesses, allergies); err != nil {
		return nil, err
	}
	return &updated, nil
}

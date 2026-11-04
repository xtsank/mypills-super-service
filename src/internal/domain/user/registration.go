package user

import (
	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

func NewRegisteredUser(id uuid.UUID, login, email, hashedPassword string, isAdmin, sex bool,
	weight, age int, isPregnant, isDriver bool, illnesses, allergies []uuid.UUID) (*User, error) {
	if len(login) < 8 || len(login) > 64 {
		return nil, errors.ErrLoginTooShort.WithSource()
	}
	if !validEmail(email) {
		return nil, errors.ErrEmailRequired.WithSource()
	}
	if !uniqueIDs(illnesses) || !uniqueIDs(allergies) {
		return nil, errors.ErrInvalidInput.WithSource()
	}
	return NewUser(id, login, email, hashedPassword, isAdmin, sex, weight, age,
		isPregnant, isDriver, illnesses, allergies)
}

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 128 {
		return errors.ErrPasswordTooShort.WithSource()
	}
	return nil
}

func ValidateLoginCredentials(login, password string) error {
	if login == "" || password == "" {
		return errors.ErrInvalidInput.WithSource()
	}
	return nil
}

package service

import (
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type IPasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type BcryptHasher struct{}

func NewBcryptHasher(i do.Injector) (IPasswordHasher, error) {
	return &BcryptHasher{}, nil
}

func (b *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.ErrInternal.WithError(err)
	}
	return string(bytes), err
}

func (b *BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

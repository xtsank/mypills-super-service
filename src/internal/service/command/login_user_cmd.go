package command

import (
	"github.com/xtsank/mypills-super-service/src/internal/errors"
)

type LoginUserCmd struct {
	Login    string
	Password string
}

func NewLoginUserCmd(
	login string,
	password string,
) (*LoginUserCmd, error) {
	if len(login) < 8 {
		return nil, errors.ErrLoginTooShort.WithSource()
	}

	if len(password) < 8 {
		return nil, errors.ErrPasswordTooShort.WithSource()
	}

	return &LoginUserCmd{
		Login:    login,
		Password: password,
	}, nil
}

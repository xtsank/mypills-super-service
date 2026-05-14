package command

type LoginUserCmd struct {
	Login    string
	Password string
}

func NewLoginUserCmd(
	login string,
	password string,
) (*LoginUserCmd, error) {
	return &LoginUserCmd{
		Login:    login,
		Password: password,
	}, nil
}

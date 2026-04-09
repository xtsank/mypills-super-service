package errors

import "net/http"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`

	HTTPStatus int
	underlying error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.underlying
}

func (e *AppError) WithError(err error) *AppError {
	e.underlying = err
	return e
}

func New(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

var (
	ErrUserExists = New("USER_EXISTS", "user with this login already exists", http.StatusConflict)

	ErrInternal = New("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
)

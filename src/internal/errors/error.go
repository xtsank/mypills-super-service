package errors

import (
	"net/http"
	"runtime"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`

	HTTPStatus int `json:"-"`
	underlying error

	SourceFile string `json:"-"`
	SourceLine int    `json:"-"`
	SourceFunc string `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.underlying
}

func (e *AppError) WithSource() *AppError {
	cloned := *e
	if cloned.SourceFile == "" {
		cloned.captureSource(1)
	}
	return &cloned
}

func (e *AppError) WithError(err error) *AppError {
	cloned := *e
	cloned.underlying = err

	if cloned.SourceFile == "" {
		cloned.captureSource(1)
	}

	return &cloned
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}

	return e.Code == t.Code
}

func (e *AppError) captureSource(skip int) {
	pc, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return
	}

	e.SourceFile = file
	e.SourceLine = line

	if fn := runtime.FuncForPC(pc); fn != nil {
		e.SourceFunc = fn.Name()
	}
}

func New(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

var (
	ErrLoginTooShort    = New("SHORT_LOGIN", "login too short", http.StatusBadRequest)
	ErrPasswordTooShort = New("SHORT_PASSWORD", "password too short", http.StatusBadRequest)
	ErrWrongWeight      = New("WRONG_WEIGHT", "weight is wrong", http.StatusBadRequest)
	ErrWrongAge         = New("WRONG_AGE", "age is wrong", http.StatusBadRequest)

	ErrQtyTooLow           = New("LOW_QTY", "quantity must be greater than zero", http.StatusBadRequest)
	ErrDateTooLate         = New("LATE_DATE", "date of manufacture cannot be in the future", http.StatusBadRequest)
	ErrCabinetItemNotFound = New("CABINET_ITEM_NOT_FOUND", "item not found", http.StatusBadRequest)

	ErrEmptyName            = New("EMPTY_NAME", "name cannot be empty", http.StatusBadRequest)
	ErrExpireTimeTooLow     = New("LOW_EXPIRE_TIME", "expire time must be greater than zero", http.StatusBadRequest)
	ErrInvalidConcentration = New("INVALID_CONCENENT", "invalid concentration", http.StatusBadRequest)
	ErrInvalidDosageRange   = New("INVALID_DOSAGE_RANGE", "dosage range is invalid", http.StatusBadRequest)
	ErrInvalidDosageValue   = New("INVALID_DOSAGE_VALUE", "dosage value is invalid", http.StatusBadRequest)
	ErrInvalidNumDoses      = New("INVALID_NUM_DOSSES", "number of doses is invalid", http.StatusBadRequest)

	ErrMedicineNotFound = New("MEDICINE_NOT_FOUND", "medicine not found", http.StatusNotFound)

	ErrInvalidCredentials = New("INVALID_CREDENTIALS", "invalid login or password", http.StatusUnauthorized)
	ErrUserNotFound       = New("USER_NOT_FOUND", "user not found", http.StatusNotFound)

	ErrUserExists = New("USER_EXISTS", "user with this login already exists", http.StatusConflict)

	ErrInternal     = New("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
	ErrInvalidInput = New("INVALID_INPUT", "invalid input data provided", http.StatusBadRequest)

	ErrUnauthorized = New("UNAUTHORIZED", "unauthorized access", http.StatusUnauthorized)
	ErrNotFound    = New("NOT_FOUND", "resource not found", http.StatusNotFound)
)

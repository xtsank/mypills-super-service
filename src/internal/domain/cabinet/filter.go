package cabinet

import apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"

func ValidateExpiryFilter(days *int, expired bool) error {
	if days != nil && (*days < 0 || *days > 365 || expired) {
		return apperrors.ErrInvalidInput.WithSource()
	}
	return nil
}

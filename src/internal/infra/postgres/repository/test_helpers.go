package repository

import "errors"

func unwrapErr(err error) string {
	if err == nil {
		return ""
	}
	if unwrapped := errors.Unwrap(err); unwrapped != nil {
		return unwrapped.Error()
	}
	return err.Error()
}


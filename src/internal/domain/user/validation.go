package user

import (
	"net/mail"
	"strings"

	"github.com/google/uuid"
)

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && strings.Contains(value, "@")
}

func uniqueIDs(values []uuid.UUID) bool {
	seen := make(map[uuid.UUID]bool, len(values))
	for _, id := range values {
		if id == uuid.Nil || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}

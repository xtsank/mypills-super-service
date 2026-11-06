package cabinet

import (
	"strings"

	"github.com/google/uuid"
	apperrors "github.com/xtsank/mypills-super-service/src/internal/errors"
)

type Cabinet struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
}

func NewCabinet(id, userID uuid.UUID, name string) (*Cabinet, error) {
	item := &Cabinet{ID: id, UserID: userID}
	if err := item.Rename(name); err != nil {
		return nil, err
	}
	return item, nil
}

func (c *Cabinet) Rename(name string) error {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > 100 {
		return apperrors.ErrInvalidInput.WithSource()
	}
	c.Name = name
	return nil
}

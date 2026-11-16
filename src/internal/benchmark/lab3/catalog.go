package lab3

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

const (
	pageSize       = 20
	medicineFields = "id, form_id, unit_id, name, expire_time, effect_on_driver, " +
		"effect_on_pregnant, method_of_application, is_prescription"
)

var escapeName = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

type catalog interface {
	Find(context.Context, uuid.UUID) (entity.MedicineEntity, error)
	Count(context.Context, string) (int64, error)
	Page(context.Context, string, int) ([]entity.MedicineEntity, error)
	Insert(context.Context, entity.MedicineEntity) error
	Rename(context.Context, uuid.UUID, string) error
	Delete(context.Context, uuid.UUID) error
}

func namePattern(name string) string {
	return "%" + escapeName.Replace(name) + "%"
}

func requireOneRow(rows int64) error {
	switch rows {
	case 1:
		return nil
	case 0:
		return sql.ErrNoRows
	default:
		return fmt.Errorf("write affected %d rows, expected one", rows)
	}
}

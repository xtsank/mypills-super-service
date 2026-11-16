package lab3

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
)

const searchName = "альфа"

type operation struct {
	name string
	run  func(context.Context, catalog, int) error
}

type scenario struct {
	name       string
	operations []operation
}

func scenarios(operations int) []scenario {
	findIDs := make([]uuid.UUID, operations)
	renameIDs := make([]uuid.UUID, operations)
	deleteIDs := make([]uuid.UUID, operations)
	names := make([]string, operations)
	inserts := make([]entity.MedicineEntity, operations)
	for i := range operations {
		findIDs[i] = medicineID(1 + i)
		renameIDs[i] = medicineID(7501 + i)
		deleteIDs[i] = medicineID(5001 + i)
		names[i] = fmt.Sprintf("Демо Бета бенчмарк %05d", i)
		inserts[i] = newMedicine(10001+i, names[i])
	}

	find := operation{"find", func(ctx context.Context, c catalog, i int) error {
		_, err := c.Find(ctx, findIDs[i])
		return err
	}}
	count := operation{"count", func(ctx context.Context, c catalog, _ int) error {
		_, err := c.Count(ctx, searchName)
		return err
	}}
	page := operation{"page", func(ctx context.Context, c catalog, i int) error {
		_, err := c.Page(ctx, searchName, (i%25)*pageSize)
		return err
	}}
	insert := operation{"insert", func(ctx context.Context, c catalog, i int) error {
		return c.Insert(ctx, inserts[i])
	}}
	rename := operation{"rename", func(ctx context.Context, c catalog, i int) error {
		return c.Rename(ctx, renameIDs[i], names[i])
	}}
	remove := operation{"delete", func(ctx context.Context, c catalog, i int) error {
		return c.Delete(ctx, deleteIDs[i])
	}}
	return []scenario{
		{find.name, []operation{find}},
		{count.name, []operation{count}},
		{page.name, []operation{page}},
		{insert.name, []operation{insert}},
		{rename.name, []operation{rename}},
		{remove.name, []operation{remove}},
		{"mixed", []operation{page, rename}},
	}
}

func medicineID(number int) uuid.UUID {
	return uuid.MustParse(fmt.Sprintf("20000000-0000-4000-8000-%012d", number))
}

func newMedicine(number int, name string) entity.MedicineEntity {
	return entity.MedicineEntity{
		ID:                  medicineID(number),
		FormId:              uuid.MustParse("10000000-0000-4000-8000-000000000001"),
		UnitId:              uuid.MustParse("10000000-0000-4000-8000-000000000002"),
		Name:                name,
		ExpireTime:          24,
		MethodOfApplication: "Тестовое описание",
	}
}

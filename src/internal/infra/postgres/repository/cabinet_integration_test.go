package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
)

func TestPostgresCabinetRepository_SaveFindUpdateDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PostgresCabinetItemRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// create user and medicine required for FK
	uid := uuid.New()
	mid := uuid.New()
	form := uuid.New()
	unit := uuid.New()

	if _, err := db.Exec("insert into Users (id, login, password, is_admin, sex, weight, age, is_pregnant, is_driver) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		uid, "u1_"+uuid.NewString(), "p", false, false, 70, 30, false, false); err != nil {
		t.Fatalf("insert user: %v", err)
	}
	if _, err := db.Exec("insert into Form (id, name) values ($1,$2)", form, "f"); err != nil {
		t.Fatalf("insert form: %v", err)
	}
	if _, err := db.Exec("insert into Unit (id, name) values ($1,$2)", unit, "u"); err != nil {
		t.Fatalf("insert unit: %v", err)
	}
	if _, err := db.Exec(`insert into Medicine (id, form_id, unit_id, name, expire_time, effect_on_driver, effect_on_pregnant, method_of_application, is_prescription) values ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		mid, form, unit, "m", 12, false, false, "oral", false); err != nil {
		t.Fatalf("insert medicine: %v", err)
	}

	cleanupCabinetItem(t, db, uid, mid, form, unit)

	item, err := cabinet_item.NewCabinetItem(uuid.New(), uid, mid, time.Now().Add(-time.Hour), 2)
	if err != nil {
		t.Fatalf("new cabinet item: %v", err)
	}
	if err := repo.Save(ctx, item); err != nil {
		t.Fatalf("save item: %s", unwrapErr(err))
	}

	items, err := repo.FindByUserID(ctx, uid)
	if err != nil {
		t.Fatalf("find by user: %v", err)
	}
	var got *cabinet_item.CabinetItem
	for _, it := range items {
		if it.ID == item.ID {
			got = it
			break
		}
	}
	if got == nil {
		t.Fatalf("expected item by id not found")
	}

	if got.Quantity != item.Quantity {
		t.Fatalf("quantity mismatch")
	}

	// update
	got.Quantity = 5
	if err := repo.Update(ctx, got); err != nil {
		t.Fatalf("update item: %s", unwrapErr(err))
	}
	after, err := repo.FindById(ctx, got.ID)
	if err != nil {
		t.Fatalf("find by id: %s", unwrapErr(err))
	}
	if after == nil || after.Quantity != 5 {
		t.Fatalf("update failed")
	}

	if err := repo.Delete(ctx, got.ID); err != nil {
		t.Fatalf("delete item: %s", unwrapErr(err))
	}
}

func cleanupCabinetItem(t *testing.T, db *sqlx.DB, userID, medicineID, formID, unitID uuid.UUID) {
	t.Helper()

	t.Cleanup(func() {
		if medicineID != uuid.Nil {
			_, _ = db.Exec("delete from User_Medicine where medicine_id = $1", medicineID)
			_, _ = db.Exec("delete from Medicine where id = $1", medicineID)
		}
		if formID != uuid.Nil {
			_, _ = db.Exec("delete from Form where id = $1", formID)
		}
		if unitID != uuid.Nil {
			_, _ = db.Exec("delete from Unit where id = $1", unitID)
		}
		if userID != uuid.Nil {
			_, _ = db.Exec("delete from Users where id = $1", userID)
		}
	})
}

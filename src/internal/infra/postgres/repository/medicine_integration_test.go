package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
)

func TestPostgresMedicineRepository_Create_AddDosage_UpdateIndications(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PostgresMedicineRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// need form and unit
	formID := uuid.New()
	unitID := uuid.New()
	if _, err := db.Exec("insert into Form (id, name) values ($1,$2)", formID, "f"); err != nil {
		t.Fatalf("insert form: %v", err)
	}
	if _, err := db.Exec("insert into Unit (id, name) values ($1,$2)", unitID, "u"); err != nil {
		t.Fatalf("insert unit: %v", err)
	}

	med, err := medicine.NewMedicine(
		uuid.New(),
		"int_med",
		6,
		false,
		"oral",
		false,
		false,
		formID,
		unitID,
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("new medicine: %v", err)
	}

	cleanupMedicine(t, db, med.ID, formID, unitID, nil)

	if err := repo.Create(ctx, med); err != nil {
		t.Fatalf("create med: %s", unwrapErr(err))
	}

	got, err := repo.FindByID(ctx, med.ID)
	if err != nil {
		t.Fatalf("find med: %s", unwrapErr(err))
	}
	if got == nil || got.Name != med.Name {
		t.Fatalf("unexpected med: %#v", got)
	}

	// add dosage rule
	rule := &medicine.DosageRule{ID: uuid.New(), ValueFrom: 1, ValueTo: 100, Type: medicine.ByWeight, DosageValue: 1.0, NumberOfDosesPerDay: 2}
	if err := repo.AddDosageRule(ctx, med.ID, rule); err != nil {
		t.Fatalf("add dosage: %s", unwrapErr(err))
	}

	got2, err := repo.FindByID(ctx, med.ID)
	if err != nil {
		t.Fatalf("find after add dosage: %s", unwrapErr(err))
	}
	found := false
	for _, d := range got2.Dosages {
		if d.ID == rule.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("added dosage not found")
	}

	// update indications: create illness and set
	ill := uuid.New()
	if _, err := db.Exec("insert into Illness (id, name) values ($1,$2)", ill, "x"); err != nil {
		t.Fatalf("insert illness: %v", err)
	}
	cleanupMedicine(t, db, med.ID, formID, unitID, []uuid.UUID{ill})

	if err := repo.UpdateIndications(ctx, med.ID, []uuid.UUID{ill}); err != nil {
		t.Fatalf("update indications: %s", unwrapErr(err))
	}
	got3, err := repo.FindByID(ctx, med.ID)
	if err != nil {
		t.Fatalf("find after indications: %s", unwrapErr(err))
	}
	if len(got3.Recommendation) == 0 || got3.Recommendation[0] != ill {
		t.Fatalf("indication not set: %#v", got3.Recommendation)
	}
}

func cleanupMedicine(t *testing.T, db *sqlx.DB, medicineID, formID, unitID uuid.UUID, illnessIDs []uuid.UUID) {
	t.Helper()

	t.Cleanup(func() {
		if medicineID != uuid.Nil {
			_, _ = db.Exec("delete from Dosage where medicine_id = $1", medicineID)
			_, _ = db.Exec("delete from Medicine_Substance where medicine_id = $1", medicineID)
			_, _ = db.Exec("delete from Recommendations where medicine_id = $1", medicineID)
			_, _ = db.Exec("delete from Contraindications where medicine_id = $1", medicineID)
			_, _ = db.Exec("delete from Medicine where id = $1", medicineID)
		}
		if formID != uuid.Nil {
			_, _ = db.Exec("delete from Form where id = $1", formID)
		}
		if unitID != uuid.Nil {
			_, _ = db.Exec("delete from Unit where id = $1", unitID)
		}
		for _, id := range illnessIDs {
			_, _ = db.Exec("delete from Illness where id = $1", id)
		}
	})
}

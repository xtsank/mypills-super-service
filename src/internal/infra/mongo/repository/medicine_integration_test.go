package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/medicine"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestMongoMedicineRepository_Create_AddDosage_UpdateIndications(t *testing.T) {
	client := setupTestMongo(t)
	defer func() { _ = client.Disconnect(context.Background()) }()

	col := client.Database("test").Collection("medicines")
	repo := &MongoMedicineRepository{collection: col}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	med, err := medicine.NewMedicine(
		uuid.New(),
		"int_med",
		6,
		false,
		"oral",
		false,
		false,
		uuid.New(),
		uuid.New(),
		nil,
		nil,
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("new medicine: %v", err)
	}

	cleanupMongoMedicine(t, col, med.ID)

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

	ill := uuid.New()
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

func cleanupMongoMedicine(t *testing.T, col *mongo.Collection, medicineID uuid.UUID) {
	t.Helper()

	t.Cleanup(func() {
		_, _ = col.DeleteOne(context.Background(), bson.M{"_id": medicineID})
	})
}

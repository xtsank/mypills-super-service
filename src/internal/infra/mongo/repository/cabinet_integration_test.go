package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	mongoEntity "github.com/xtsank/mypills-super-service/src/internal/infra/mongo/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestMongoCabinetRepository_SaveFindUpdateDelete(t *testing.T) {
	client := setupTestMongo(t)
	defer func() { _ = client.Disconnect(context.Background()) }()

	db := client.Database("test")
	repo := &MongoCabinetItemRepository{collection: db.Collection("cabinet_items"), medicines: db.Collection("medicines")}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	uid := uuid.New()
	mid := uuid.New()

	med := mongoEntity.MedicineEntity{
		ID:         mid,
		Name:       "m",
		ExpireTime: 12,
	}
	if _, err := repo.medicines.InsertOne(ctx, med); err != nil {
		t.Fatalf("insert medicine: %v", err)
	}

	cleanupMongoCabinetItem(t, db, uid, mid)

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

func cleanupMongoCabinetItem(t *testing.T, db *mongo.Database, userID, medicineID uuid.UUID) {
	t.Helper()

	t.Cleanup(func() {
		if medicineID != uuid.Nil {
			_, _ = db.Collection("cabinet_items").DeleteMany(context.Background(), bson.M{"medicine_id": medicineID})
			_, _ = db.Collection("medicines").DeleteOne(context.Background(), bson.M{"_id": medicineID})
		}
		if userID != uuid.Nil {
			_, _ = db.Collection("cabinet_items").DeleteMany(context.Background(), bson.M{"user_id": userID})
		}
	})
}

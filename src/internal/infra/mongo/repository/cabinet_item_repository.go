package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/config"
	"github.com/xtsank/mypills-super-service/src/internal/domain/cabinet_item"
	"github.com/xtsank/mypills-super-service/src/internal/dto"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	mongoEntity "github.com/xtsank/mypills-super-service/src/internal/infra/mongo/entity"
	postgresEntity "github.com/xtsank/mypills-super-service/src/internal/infra/postgres/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoCabinetItemRepository struct {
	collection *mongo.Collection
	medicines  *mongo.Collection
}

func NewMongoCabinetItemRepository(i do.Injector) (cabinet_item.ICabinetItemRepository, error) {
	client := do.MustInvoke[*mongo.Client](i)
	cfg := do.MustInvoke[*config.Config](i)

	db := client.Database(cfg.DBName)
	return &MongoCabinetItemRepository{
		collection: db.Collection("cabinet_items"),
		medicines:  db.Collection("medicines"),
	}, nil
}

func (r *MongoCabinetItemRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*cabinet_item.CabinetItem, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	defer cursor.Close(ctx)

	var ents []mongoEntity.CabinetItemEntity
	if err := cursor.All(ctx, &ents); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	res := make([]*cabinet_item.CabinetItem, 0, len(ents))
	for _, ent := range ents {
		item, err := cabinet_item.NewCabinetItem(ent.ID, ent.UserID, ent.MedicineID, ent.DateOfManufacture, ent.Quantity)
		if err != nil {
			return nil, err
		}
		res = append(res, item)
	}

	return res, nil
}

func (r *MongoCabinetItemRepository) FindDetailsByUserID(ctx context.Context, userID uuid.UUID) ([]postgresEntity.CabinetItemWithMedicineEntity, error) {
	items, err := r.findCabinetItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []postgresEntity.CabinetItemWithMedicineEntity{}, nil
	}

	medMap, err := r.loadMedicines(ctx, items)
	if err != nil {
		return nil, err
	}

	res := make([]postgresEntity.CabinetItemWithMedicineEntity, 0, len(items))
	for _, item := range items {
		med, ok := medMap[item.MedicineID]
		if !ok {
			continue
		}
		expiresAt := item.DateOfManufacture.AddDate(0, med.ExpireTime, 0)
		res = append(res, postgresEntity.CabinetItemWithMedicineEntity{
			ID:                item.ID,
			UserID:            item.UserID,
			MedicineID:        item.MedicineID,
			MedicineName:      med.Name,
			DateOfManufacture: item.DateOfManufacture,
			Quantity:          item.Quantity,
			ExpiresAt:         expiresAt,
		})
	}

	return res, nil
}

func (r *MongoCabinetItemRepository) FindExistingCabinetItem(ctx context.Context, userID uuid.UUID, medID uuid.UUID, date time.Time) (*cabinet_item.CabinetItem, error) {
	var ent mongoEntity.CabinetItemEntity
	filter := bson.M{
		"user_id":             userID,
		"medicine_id":         medID,
		"date_of_manufacture": date,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&ent)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return cabinet_item.NewCabinetItem(ent.ID, ent.UserID, ent.MedicineID, ent.DateOfManufacture, ent.Quantity)
}

func (r *MongoCabinetItemRepository) FindById(ctx context.Context, id uuid.UUID) (*cabinet_item.CabinetItem, error) {
	var ent mongoEntity.CabinetItemEntity
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&ent)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	return cabinet_item.NewCabinetItem(ent.ID, ent.UserID, ent.MedicineID, ent.DateOfManufacture, ent.Quantity)
}

func (r *MongoCabinetItemRepository) FindExpiredByUserID(ctx context.Context, userID uuid.UUID) ([]*dto.ExpiredItemDto, error) {
	items, err := r.FindDetailsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	res := make([]*dto.ExpiredItemDto, 0)
	for _, item := range items {
		if item.ExpiresAt.Before(today) {
			res = append(res, &dto.ExpiredItemDto{
				ID:                item.ID,
				UserID:            item.UserID,
				MedicineID:        item.MedicineID,
				MedicineName:      item.MedicineName,
				DateOfManufacture: item.DateOfManufacture,
				Quantity:          item.Quantity,
				ExpiresAt:         item.ExpiresAt,
			})
		}
	}

	return res, nil
}

func (r *MongoCabinetItemRepository) Update(ctx context.Context, item *cabinet_item.CabinetItem) error {
	filter := bson.M{"_id": item.ID}
	update := bson.M{"$set": bson.M{"quantity": item.Quantity}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *MongoCabinetItemRepository) Save(ctx context.Context, item *cabinet_item.CabinetItem) error {
	ent := mongoEntity.CabinetItemEntity{
		ID:                item.ID,
		UserID:            item.UserID,
		MedicineID:        item.MedicineID,
		DateOfManufacture: item.DateOfManufacture,
		Quantity:          item.Quantity,
	}

	if _, err := r.collection.InsertOne(ctx, ent); err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *MongoCabinetItemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *MongoCabinetItemRepository) findCabinetItems(ctx context.Context, userID uuid.UUID) ([]mongoEntity.CabinetItemEntity, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	defer cursor.Close(ctx)

	var items []mongoEntity.CabinetItemEntity
	if err := cursor.All(ctx, &items); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	return items, nil
}

func (r *MongoCabinetItemRepository) loadMedicines(ctx context.Context, items []mongoEntity.CabinetItemEntity) (map[uuid.UUID]mongoEntity.MedicineEntity, error) {
	ids := make([]uuid.UUID, 0, len(items))
	idSet := map[uuid.UUID]struct{}{}
	for _, item := range items {
		if _, ok := idSet[item.MedicineID]; ok {
			continue
		}
		idSet[item.MedicineID] = struct{}{}
		ids = append(ids, item.MedicineID)
	}

	opts := options.Find().SetProjection(bson.M{"name": 1, "expire_time": 1})
	cursor, err := r.medicines.Find(ctx, bson.M{"_id": bson.M{"$in": ids}}, opts)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	defer cursor.Close(ctx)

	var meds []mongoEntity.MedicineEntity
	if err := cursor.All(ctx, &meds); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	res := make(map[uuid.UUID]mongoEntity.MedicineEntity, len(meds))
	for _, med := range meds {
		res[med.ID] = med
	}
	return res, nil
}

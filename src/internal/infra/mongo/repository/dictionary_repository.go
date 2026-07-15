package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/config"
	"github.com/xtsank/mypills-super-service/src/internal/domain/dictionary"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	mongoEntity "github.com/xtsank/mypills-super-service/src/internal/infra/mongo/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDictionaryRepository struct {
	db *mongo.Database
}

func NewMongoDictionaryRepository(i do.Injector) (dictionary.IRepository, error) {
	client := do.MustInvoke[*mongo.Client](i)
	cfg := do.MustInvoke[*config.Config](i)

	return &MongoDictionaryRepository{db: client.Database(cfg.DBName)}, nil
}

func (r *MongoDictionaryRepository) ListIllnesses(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, "illnesses")
}

func (r *MongoDictionaryRepository) ListSubstances(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, "substances")
}

func (r *MongoDictionaryRepository) ListForms(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, "forms")
}

func (r *MongoDictionaryRepository) ListUnits(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, "units")
}

func (r *MongoDictionaryRepository) ListMedicines(ctx context.Context) ([]dictionary.Item, error) {
	return r.list(ctx, "medicines")
}

func (r *MongoDictionaryRepository) list(ctx context.Context, collection string) ([]dictionary.Item, error) {
	opts := options.Find().SetProjection(bson.M{"name": 1}).SetSort(bson.M{"name": 1})
	cursor, err := r.db.Collection(collection).Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	defer cursor.Close(ctx)

	if collection == "medicines" {
		var meds []struct {
			ID   uuid.UUID `bson:"_id"`
			Name string    `bson:"name"`
		}
		if err := cursor.All(ctx, &meds); err != nil {
			return nil, appErrors.ErrInternal.WithError(err)
		}
		res := make([]dictionary.Item, 0, len(meds))
		for _, med := range meds {
			res = append(res, dictionary.Item{ID: med.ID, Name: med.Name})
		}
		return res, nil
	}

	var ents []mongoEntity.DictionaryItemEntity
	if err := cursor.All(ctx, &ents); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	items := make([]dictionary.Item, 0, len(ents))
	for _, ent := range ents {
		items = append(items, dictionary.Item{ID: ent.ID, Name: ent.Name})
	}
	return items, nil
}

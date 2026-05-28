package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
	"github.com/xtsank/mypills-super-service/src/internal/config"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	appErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/infra/mongo/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(i do.Injector) (user.IUserRepository, error) {
	client := do.MustInvoke[*mongo.Client](i)
	cfg := do.MustInvoke[*config.Config](i)

	collection := client.Database(cfg.DBName).Collection("users")

	return &MongoUserRepository{collection: collection}, nil
}

func (r *MongoUserRepository) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	var ent entity.UserEntity
	filter := bson.M{"login": login}
	err := r.collection.FindOne(ctx, filter).Decode(&ent)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrUserNotFound.WithSource()
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	notify := &user.NotifyInfo{
		Enabled:         ent.NotifyEnabled,
		IntervalMinutes: ent.NotifyIntervalMinutes,
		LastNotifiedAt:  ent.LastNotifiedAt,
	}

	return user.NewUser(
		ent.ID,
		ent.Login,
		ent.Email,
		ent.Password,
		ent.IsAdmin,
		ent.Sex,
		ent.Weight,
		ent.Age,
		ent.IsPregnant,
		ent.IsDriver,
		notify,
		ent.Illnesses,
		ent.Allergies,
	)
}

func (r *MongoUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var ent entity.UserEntity
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&ent)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErrors.ErrUserNotFound.WithSource()
		}
		return nil, appErrors.ErrInternal.WithError(err)
	}

	notify := &user.NotifyInfo{
		Enabled:         ent.NotifyEnabled,
		IntervalMinutes: ent.NotifyIntervalMinutes,
		LastNotifiedAt:  ent.LastNotifiedAt,
	}

	return user.NewUser(
		ent.ID,
		ent.Login,
		ent.Email,
		ent.Password,
		ent.IsAdmin,
		ent.Sex,
		ent.Weight,
		ent.Age,
		ent.IsPregnant,
		ent.IsDriver,
		notify,
		ent.Illnesses,
		ent.Allergies,
	)
}

func (r *MongoUserRepository) FindNotifyEnabled(ctx context.Context) ([]*user.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"notify_enabled": true})
	if err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}
	defer cursor.Close(ctx)

	var ents []entity.UserEntity
	if err := cursor.All(ctx, &ents); err != nil {
		return nil, appErrors.ErrInternal.WithError(err)
	}

	users := make([]*user.User, 0, len(ents))
	for _, ent := range ents {
		notify := &user.NotifyInfo{
			Enabled:         ent.NotifyEnabled,
			IntervalMinutes: ent.NotifyIntervalMinutes,
			LastNotifiedAt:  ent.LastNotifiedAt,
		}
		u, err := user.NewUser(
			ent.ID,
			ent.Login,
			ent.Email,
			ent.Password,
			ent.IsAdmin,
			ent.Sex,
			ent.Weight,
			ent.Age,
			ent.IsPregnant,
			ent.IsDriver,
			notify,
			nil,
			nil,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *MongoUserRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	filter := bson.M{"login": login}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, appErrors.ErrInternal.WithError(err)
	}
	return count > 0, nil
}

func (r *MongoUserRepository) Create(ctx context.Context, u *user.User) error {
	ent := entity.UserEntity{
		ID:                    u.ID,
		Login:                 u.Login,
		Email:                 u.Email,
		Password:              u.Password,
		IsAdmin:               u.IsAdmin,
		Sex:                   u.Sex,
		Weight:                u.Weight,
		Age:                   u.Age,
		IsPregnant:            u.IsPregnant,
		IsDriver:              u.IsDriver,
		NotifyEnabled:         u.Notify.Enabled,
		NotifyIntervalMinutes: u.Notify.IntervalMinutes,
		LastNotifiedAt:        u.Notify.LastNotifiedAt,
		Illnesses:             u.Illnesses,
		Allergies:             u.Allergies,
	}

	_, err := r.collection.InsertOne(ctx, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}
	return nil
}

func (r *MongoUserRepository) Update(ctx context.Context, u *user.User) error {
	ent := entity.UserEntity{
		ID:                    u.ID,
		Login:                 u.Login,
		Email:                 u.Email,
		Password:              u.Password,
		IsAdmin:               u.IsAdmin,
		Sex:                   u.Sex,
		Weight:                u.Weight,
		Age:                   u.Age,
		IsPregnant:            u.IsPregnant,
		IsDriver:              u.IsDriver,
		NotifyEnabled:         u.Notify.Enabled,
		NotifyIntervalMinutes: u.Notify.IntervalMinutes,
		LastNotifiedAt:        u.Notify.LastNotifiedAt,
		Illnesses:             u.Illnesses,
		Allergies:             u.Allergies,
	}

	filter := bson.M{"_id": u.ID}

	result, err := r.collection.ReplaceOne(ctx, filter, ent)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}

	if result.MatchedCount == 0 {
		return appErrors.ErrUserNotFound.WithSource()
	}

	return nil
}

func (r *MongoUserRepository) UpdateNotify(ctx context.Context, id uuid.UUID, enabled bool, intervalMinutes int, lastNotifiedAt *time.Time) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"notify_enabled":          enabled,
			"notify_interval_minutes": intervalMinutes,
			"last_notified_at":        lastNotifiedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return appErrors.ErrInternal.WithError(err)
	}

	if result.MatchedCount == 0 {
		return appErrors.ErrUserNotFound.WithSource()
	}

	return nil
}

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestMongoUserRepository_CreateAndFind(t *testing.T) {
	client := setupTestMongo(t)
	defer func() { _ = client.Disconnect(context.Background()) }()

	repo := &MongoUserRepository{collection: client.Database("test").Collection("users")}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := uuid.New()
	illnessID := uuid.New()

	u, err := user.NewUser(
		userID,
		"test_pilot_"+uuid.NewString(),
		"test@example.com",
		"hash_string",
		false,
		true,
		75,
		25,
		false,
		true,
		nil,
		[]uuid.UUID{illnessID},
		nil,
	)
	assert.NoError(t, err)

	cleanupMongoUser(t, repo.collection, userID)

	err = repo.Create(ctx, u)
	assert.NoError(t, err)

	found, err := repo.FindByLogin(ctx, u.Login)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	if found == nil {
		return
	}
	assert.Equal(t, u.ID, found.ID)
	assert.Len(t, found.Illnesses, 1)
	assert.Equal(t, illnessID, found.Illnesses[0])
}

func cleanupMongoUser(t *testing.T, col *mongo.Collection, userID uuid.UUID) {
	t.Helper()

	t.Cleanup(func() {
		_, _ = col.DeleteOne(context.Background(), bson.M{"_id": userID})
	})
}

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

func TestPostgresUserRepository_CreateAndFind(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PostgresUserRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userID := uuid.New()
	illnessID := uuid.New()

	_, err := db.Exec("INSERT INTO Illness (id, name) VALUES ($1, $2)", illnessID, "Flu")
	assert.NoError(t, err)

	u, err := user.NewUser(
		userID,
		"test_pilot_"+uuid.NewString(),
		"hash_string",
		false,
		true,
		75,
		25,
		false,
		true,
		[]uuid.UUID{illnessID},
		nil,
	)
	assert.NoError(t, err)

	t.Cleanup(func() {
		_, _ = db.Exec("delete from User_Illness where user_id = $1", userID)
		_, _ = db.Exec("delete from Users where id = $1", userID)
		_, _ = db.Exec("delete from Illness where id = $1", illnessID)
	})

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

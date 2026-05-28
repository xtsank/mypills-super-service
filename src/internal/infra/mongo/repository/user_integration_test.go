package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

func TestMongoUserRepository_CreateFindUpdate(t *testing.T) {
	client := setupTestMongo(t)
	defer func() { _ = client.Disconnect(context.Background()) }()

	col := client.Database("test").Collection("users")
	repo := &MongoUserRepository{collection: col}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	login := "int_user_" + uuid.NewString()
	u, err := user.NewUser(
		uuid.New(),
		login,
		"test@example.com",
		"pwd",
		false,
		false,
		70,
		30,
		false,
		false,
		nil,
		nil,
		[]uuid.UUID{uuid.New()},
	)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}

	cleanupMongoUser(t, col, u.ID)

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create user: %v", err)
	}

	got, err := repo.FindByLogin(ctx, u.Login)
	if err != nil {
		t.Fatalf("find by login: %v", err)
	}
	if got == nil || got.ID != u.ID {
		t.Fatalf("unexpected user returned: %#v", got)
	}

	exists, err := repo.ExistsByLogin(ctx, u.Login)
	if err != nil {
		t.Fatalf("exists by login: %v", err)
	}
	if !exists {
		t.Fatalf("expected exists true")
	}

	got2, err := repo.FindByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("find by id: %v", err)
	}
	if got2 == nil || got2.Login != u.Login {
		t.Fatalf("unexpected user by id: %#v", got2)
	}

	got2.Weight = 80
	got2.Allergies = []uuid.UUID{uuid.New()}
	if err := repo.Update(ctx, got2); err != nil {
		t.Fatalf("update user: %v", err)
	}

	after, err := repo.FindByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("find after update: %v", err)
	}
	if after.Weight != 80 {
		t.Fatalf("weight not updated, got %d", after.Weight)
	}
}

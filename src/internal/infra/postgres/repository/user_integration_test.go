package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
)

// use setupTestDB(t) from setup_test.go

func TestPostgresUserRepository_CreateFindUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := &PostgresUserRepository{db: db}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	substanceID := uuid.New()
	var newSubstanceID uuid.UUID

	if _, err := db.Exec("insert into Substance (id, name) values ($1,$2)", substanceID, "s1"); err != nil {
		t.Fatalf("insert substance: %v", err)
	}

	login := "int_user_" + uuid.NewString()
	u, err := user.NewUser(
		uuid.New(),
		login,
		"pwd",
		false,
		false,
		70,
		30,
		false,
		false,
		nil,
		[]uuid.UUID{substanceID},
	)
	if err != nil {
		t.Fatalf("new user: %v", err)
	}

	cleanupUser(t, db, u.ID, []uuid.UUID{substanceID}, &newSubstanceID)

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
	if !containsUUID(got.Allergies, substanceID) {
		t.Fatalf("allergy not set: %#v", got.Allergies)
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

	// update
	newSubstanceID = uuid.New()
	if _, err := db.Exec("insert into Substance (id, name) values ($1,$2)", newSubstanceID, "s2"); err != nil {
		t.Fatalf("insert substance 2: %v", err)
	}
	got2.Weight = 80
	got2.Allergies = []uuid.UUID{newSubstanceID}
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
	if !containsUUID(after.Allergies, newSubstanceID) {
		t.Fatalf("allergy not updated: %#v", after.Allergies)
	}
}

func cleanupUser(t *testing.T, db *sqlx.DB, userID uuid.UUID, substanceIDs []uuid.UUID, newSubstanceID *uuid.UUID) {
	t.Helper()

	t.Cleanup(func() {
		if userID != uuid.Nil {
			_, _ = db.Exec("delete from User_Substance where user_id = $1", userID)
			_, _ = db.Exec("delete from User_Illness where user_id = $1", userID)
			_, _ = db.Exec("delete from Users where id = $1", userID)
		}
		for _, sid := range substanceIDs {
			_, _ = db.Exec("delete from Substance where id = $1", sid)
		}
		if newSubstanceID != nil && *newSubstanceID != uuid.Nil {
			_, _ = db.Exec("delete from Substance where id = $1", *newSubstanceID)
		}
	})
}

func containsUUID(list []uuid.UUID, target uuid.UUID) bool {
	for _, id := range list {
		if id == target {
			return true
		}
	}
	return false
}

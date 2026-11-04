package user

import (
	"testing"

	"github.com/google/uuid"
)

func TestUpdateProfileRejectsInvalidChangeWithoutChangingUser(t *testing.T) {
	u, err := NewRegisteredUser(uuid.New(), "user0001", "old@example.test", "hashed", false,
		false, 70, 30, false, false, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	email := "new@example.test"
	age := 121
	if err := u.UpdateProfile(&email, nil, nil, &age, nil, nil, nil, nil); err == nil {
		t.Fatal("invalid age was accepted")
	}
	if u.Email != "old@example.test" || u.Age != 30 {
		t.Fatalf("profile changed after rejected update: email=%q age=%d", u.Email, u.Age)
	}
}

func TestUpdateProfileClearsAllergiesWhenEmptyListIsProvided(t *testing.T) {
	allergy := uuid.New()
	u, err := NewRegisteredUser(uuid.New(), "user0001", "user@example.test", "hashed", false,
		false, 70, 30, false, false, nil, []uuid.UUID{allergy})
	if err != nil {
		t.Fatal(err)
	}
	if err := u.UpdateProfile(nil, nil, nil, nil, nil, nil, nil, []uuid.UUID{}); err != nil {
		t.Fatal(err)
	}
	if u.Allergies == nil || len(u.Allergies) != 0 {
		t.Fatalf("expected an explicitly empty allergy list, got %v", u.Allergies)
	}
}

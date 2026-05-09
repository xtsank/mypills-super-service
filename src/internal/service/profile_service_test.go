package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/xtsank/mypills-super-service/src/internal/domain/user"
	svcErrors "github.com/xtsank/mypills-super-service/src/internal/errors"
	"github.com/xtsank/mypills-super-service/src/internal/service/command"
)

type mockUserRepoForProfile struct {
	u         *user.User
	updateErr error
}

func (m *mockUserRepoForProfile) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	return m.u, nil
}
func (m *mockUserRepoForProfile) Update(ctx context.Context, u *user.User) error { return m.updateErr }

// unused
func (m *mockUserRepoForProfile) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	panic("not used")
}
func (m *mockUserRepoForProfile) Create(ctx context.Context, u *user.User) error { panic("not used") }
func (m *mockUserRepoForProfile) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	panic("not used")
}

func TestProfileService_Update_Success(t *testing.T) {
	u := &user.User{ID: uuid.New(), Sex: false, Weight: 1}
	repo := &mockUserRepoForProfile{u: u}
	s := &ProfileService{userRepo: repo}

	age := 30
	weight := 2
	cmd := &command.UpdateProfileCmd{ID: u.ID, Age: &age, Weight: &weight}

	dto, err := s.UpdateProfile(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto, got nil")
	}
	if dto.Age != 30 {
		t.Fatalf("age not updated")
	}
}

func TestProfileService_Update_UserNotFound(t *testing.T) {
	repo := &mockUserRepoForProfile{u: nil}
	s := &ProfileService{userRepo: repo}

	_, err := s.UpdateProfile(context.Background(), &command.UpdateProfileCmd{ID: uuid.New()})
	if !errors.Is(err, svcErrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestProfileService_Update_UpdateError(t *testing.T) {
	u := &user.User{ID: uuid.New()}
	repo := &mockUserRepoForProfile{u: u, updateErr: errors.New("boom")}
	s := &ProfileService{userRepo: repo}

	_, err := s.UpdateProfile(context.Background(), &command.UpdateProfileCmd{ID: u.ID})
	if err == nil {
		t.Fatalf("expected error")
	}
}

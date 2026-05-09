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

type mockUserRepoForAuth struct {
	exists    bool
	u         *user.User
	createErr error
}

func (m *mockUserRepoForAuth) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	return m.exists, nil
}
func (m *mockUserRepoForAuth) Create(ctx context.Context, u *user.User) error { return m.createErr }
func (m *mockUserRepoForAuth) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	return m.u, nil
}

// unused methods required by interface - provide panics to catch accidental use
func (m *mockUserRepoForAuth) FindByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	panic("not implemented")
}
func (m *mockUserRepoForAuth) Update(ctx context.Context, u *user.User) error {
	panic("not implemented")
}

type mockHasher struct {
	hashVal string
	hashErr error
	cmpErr  error
}

func (m *mockHasher) Hash(password string) (string, error) { return m.hashVal, m.hashErr }
func (m *mockHasher) Compare(hash, password string) error  { return m.cmpErr }

type mockTokenManager struct {
	token  string
	genErr error
}

func (m *mockTokenManager) GenerateToken(userID uuid.UUID, isAdmin bool) (string, error) {
	return m.token, m.genErr
}
func (m *mockTokenManager) VerifyToken(tokenStr string) (uuid.UUID, bool, error) { panic("not needed") }

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := &mockUserRepoForAuth{exists: false}
	hasher := &mockHasher{hashVal: "hashed"}
	tokenManager := &mockTokenManager{token: "ttt"}

	s := &AuthService{userRepo: userRepo, hasher: hasher, tokenManager: tokenManager}

	cmd := &command.CreateUserCmd{Login: "a", Password: "p", IsAdmin: false, Weight: 70, Age: 30}

	dto, err := s.Register(context.Background(), cmd)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto == nil {
		t.Fatalf("expected dto, got nil")
	}
	if dto.Token == "" {
		t.Fatalf("expected non-empty token")
	}
}

func TestAuthService_Register_UserExists(t *testing.T) {
	userRepo := &mockUserRepoForAuth{exists: true}
	s := &AuthService{userRepo: userRepo, hasher: &mockHasher{}, tokenManager: &mockTokenManager{}}

	_, err := s.Register(context.Background(), &command.CreateUserCmd{Login: "a", Password: "p"})
	if !errors.Is(err, svcErrors.ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}

func TestAuthService_Login_NotFound(t *testing.T) {
	userRepo := &mockUserRepoForAuth{u: nil}
	s := &AuthService{userRepo: userRepo, hasher: &mockHasher{}, tokenManager: &mockTokenManager{}}

	_, err := s.Login(context.Background(), &command.LoginUserCmd{Login: "l", Password: "p"})
	if !errors.Is(err, svcErrors.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	u := &user.User{ID: uuid.New(), Login: "l", Password: "h"}
	userRepo := &mockUserRepoForAuth{u: u}
	s := &AuthService{userRepo: userRepo, hasher: &mockHasher{cmpErr: errors.New("bad")}, tokenManager: &mockTokenManager{}}

	_, err := s.Login(context.Background(), &command.LoginUserCmd{Login: "l", Password: "p"})
	if !errors.Is(err, svcErrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	u := &user.User{ID: uuid.New(), Login: "l", Password: "h", IsAdmin: true}
	userRepo := &mockUserRepoForAuth{u: u}
	s := &AuthService{userRepo: userRepo, hasher: &mockHasher{}, tokenManager: &mockTokenManager{token: "tok"}}

	dto, err := s.Login(context.Background(), &command.LoginUserCmd{Login: "l", Password: "p"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if dto.Token == "" {
		t.Fatalf("expected non-empty token")
	}
}

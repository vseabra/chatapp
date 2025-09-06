package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// mockRepository implements Repository for unit tests
type mockRepository struct {
	createFn      func(ctx context.Context, u *User) error
	findByEmailFn func(ctx context.Context, email string) (*User, error)
}

func (m *mockRepository) Create(ctx context.Context, u *User) error {
	if m.createFn != nil {
		return m.createFn(ctx, u)
	}
	return nil
}

func (m *mockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	if m.findByEmailFn != nil {
		return m.findByEmailFn(ctx, email)
	}
	return nil, mongo.ErrNoDocuments
}

func TestService_Register_Success(t *testing.T) {
	repo := &mockRepository{
		createFn: func(ctx context.Context, u *User) error {
			u.ID = primitive.NewObjectID()
			return nil
		},
	}
	s := NewService(repo)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := s.Register(ctx, RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.ID == "" || resp.Email != "alice@example.com" || resp.Name != "Alice" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestService_Register_DuplicateEmail(t *testing.T) {
	repo := &mockRepository{
		createFn: func(ctx context.Context, u *User) error {
			return mongo.CommandError{Code: 11000}
		},
	}
	s := NewService(repo)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := s.Register(ctx, RegisterRequest{Name: "Bob", Email: "bob@example.com", Password: "password"})
	if !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("expected ErrEmailTaken, got %v", err)
	}
}

func TestService_Register_InvalidInput(t *testing.T) {
	s := NewService(&mockRepository{})
	ctx := context.Background()
	_, err := s.Register(ctx, RegisterRequest{})
	if err == nil {
		t.Fatalf("expected error for invalid input")
	}
}

func TestService_Login_Success(t *testing.T) {
	hashed, _ := bcryptGenerate("secret123")
	repo := &mockRepository{
		findByEmailFn: func(ctx context.Context, email string) (*User, error) {
			return &User{Email: email, Password: hashed}, nil
		},
	}
	s := NewService(repo)
	ctx := context.Background()
	resp, err := s.Login(ctx, LoginRequest{Email: "alice@example.com", Password: "secret123"}, "testsecret", "1h")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp == nil || resp.AccessToken == "" {
		t.Fatalf("expected access token, got %+v", resp)
	}
}

func TestService_Login_InvalidCredentials(t *testing.T) {
	hashed, _ := bcryptGenerate("rightpass")
	repo := &mockRepository{
		findByEmailFn: func(ctx context.Context, email string) (*User, error) {
			return &User{Email: email, Password: hashed}, nil
		},
	}
	s := NewService(repo)
	ctx := context.Background()
	_, err := s.Login(ctx, LoginRequest{Email: "x@example.com", Password: "wrong"}, "secret", "1h")
	if err == nil {
		t.Fatalf("expected error for invalid credentials")
	}
}

func TestService_Login_UserNotFound(t *testing.T) {
	repo := &mockRepository{findByEmailFn: func(ctx context.Context, email string) (*User, error) {
		return nil, mongo.ErrNoDocuments
	}}
	s := NewService(repo)
	ctx := context.Background()
	_, err := s.Login(ctx, LoginRequest{Email: "missing@example.com", Password: "pass"}, "secret", "1h")
	if err == nil {
		t.Fatalf("expected error when user not found")
	}
}

// helper to avoid importing bcrypt in test imports section of each test
func bcryptGenerate(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

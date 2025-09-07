package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"chatapp/internal/auth"
)

var ErrEmailTaken = errors.New("email already registered")

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest, jwtSecret string, expiresIn string) (*LoginResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if name == "" || email == "" || password == "" {
		return nil, errors.New("invalid input")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:      name,
		Email:     email,
		Password:  string(hash),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		// Detect duplicate key error (code 11000)
		var mw mongo.WriteException
		if errors.As(err, &mw) {
			for _, we := range mw.WriteErrors {
				if we.Code == 11000 {
					return nil, ErrEmailTaken
				}
			}
		}
		var ce mongo.CommandError
		if errors.As(err, &ce) && ce.Code == 11000 {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	resp := &RegisterResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}
	return resp, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest, jwtSecret string, expiresIn string) (*LoginResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password
	if email == "" || password == "" {
		return nil, errors.New("invalid credentials")
	}
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	token, err := auth.GenerateToken(jwtSecret, expiresIn, u.ID.Hex(), u.Email)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{AccessToken: token, UserId: u.ID.Hex(), UserName: u.Name}, nil
}

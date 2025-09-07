package message

import (
	"context"
	"errors"
	"strings"
	"time"
)

type ChatRoomReader interface {
	// Minimal method used for membership/room existence checks.
	FindByID(ctx context.Context, id string) (bool, error)
}

type Service interface {
	Create(ctx context.Context, userID string, roomID string, text string) (*Message, error)
	List(ctx context.Context, roomID string, limit int64, cursor string) ([]Message, string, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, userID string, roomID string, text string) (*Message, error) {
	t := strings.TrimSpace(text)
	if t == "" {
		return nil, errors.New("empty message")
	}
	m := &Message{
		RoomID:    roomID,
		UserID:    userID,
		Text:      t,
		Type:      "user",
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Insert(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *service) List(ctx context.Context, roomID string, limit int64, cursor string) ([]Message, string, error) {
	return s.repo.ListByRoom(ctx, roomID, limit, cursor)
}

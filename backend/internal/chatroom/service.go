package chatroom

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrNotFound = errors.New("chatroom not found")
var ErrForbidden = errors.New("forbidden")

type Service interface {
	Create(ctx context.Context, ownerID string, req CreateChatRoomRequest) (*ChatRoomResponse, error)
	ListAll(ctx context.Context, limit int64, skip int64) ([]ChatRoomResponse, error)
	Rename(ctx context.Context, userID string, id string, req UpdateChatRoomRequest) error
	Delete(ctx context.Context, userID string, id string) error
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(ctx context.Context, ownerID string, req CreateChatRoomRequest) (*ChatRoomResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.New("invalid input")
	}
	room := &ChatRoom{
		Title:     title,
		OwnerID:   ownerID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, room); err != nil {
		return nil, err
	}
	return &ChatRoomResponse{ID: room.ID.Hex(), Title: room.Title, OwnerID: room.OwnerID}, nil
}

func (s *service) ListAll(ctx context.Context, limit int64, skip int64) ([]ChatRoomResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if skip < 0 {
		skip = 0
	}
	items, err := s.repo.FindAll(ctx, limit, skip)
	if err != nil {
		return nil, err
	}
	resp := make([]ChatRoomResponse, 0, len(items))
	for _, c := range items {
		resp = append(resp, ChatRoomResponse{ID: c.ID.Hex(), Title: c.Title, OwnerID: c.OwnerID})
	}
	return resp, nil
}

func (s *service) Rename(ctx context.Context, userID string, id string, req UpdateChatRoomRequest) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}
	c, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return ErrNotFound
	}
	if c.OwnerID != userID {
		return ErrForbidden
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return errors.New("invalid input")
	}
	return s.repo.UpdateTitle(ctx, oid, title)
}

func (s *service) Delete(ctx context.Context, userID string, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrNotFound
	}
	ok, err := s.repo.Delete(ctx, oid, userID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

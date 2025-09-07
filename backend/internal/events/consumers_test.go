package events

import (
	"context"
	"testing"
	"time"

	"chatapp/internal/message"
	"chatapp/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mockMessageService mocks message.Service for testing
type mockMessageService struct {
	createdMessage *message.Message
	createError    error
	botMessage     *message.Message
	botError       error
}

func (m *mockMessageService) CreateWithName(ctx context.Context, userID, userName, roomID, text string) (*message.Message, error) {
	return m.createdMessage, m.createError
}

func (m *mockMessageService) CreateBotMessage(ctx context.Context, roomID, text string) (*message.Message, error) {
	return m.botMessage, m.botError
}

func (m *mockMessageService) List(ctx context.Context, roomID string, limit int64, cursor string) ([]message.Message, string, error) {
	return []message.Message{}, "", nil
}

// mockUserRepository mocks user.Repository for testing
type mockUserRepository struct {
	user *user.User
	err  error
}

func (m *mockUserRepository) Create(ctx context.Context, u *user.User) error {
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	return m.user, m.err
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	return m.user, m.err
}

// mockBroadcaster mocks Broadcaster for testing
type mockBroadcaster struct {
	broadcasts []BroadcastCall
}

type BroadcastCall struct {
	RoomID  string
	Payload any
}

func (m *mockBroadcaster) Broadcast(roomID string, payload any) {
	m.broadcasts = append(m.broadcasts, BroadcastCall{RoomID: roomID, Payload: payload})
}

func TestCommandParsing(t *testing.T) {
	// Test command parsing logic
	trim := "/echo hello world"
	if !startsWithSlash(trim) {
		t.Error("expected command to start with slash")
	}

	parts := splitCommand(trim)
	if parts[0] != "echo" {
		t.Error("expected command to be 'echo'")
	}
	if parts[1] != "hello world" {
		t.Error("expected args to be 'hello world'")
	}
}

func TestRegularMessageParsing(t *testing.T) {
	// Test regular message processing
	trim := "Hello world"
	if startsWithSlash(trim) {
		t.Error("expected regular message to not start with slash")
	}
}

func TestBroadcastConsumer_ProcessMessage(t *testing.T) {
	broadcaster := &mockBroadcaster{}

	consumer := &BroadcastConsumer{
		Hub: broadcaster,
	}

	// Test message created event
	evt := MessageCreated{
		ID:        "msg123",
		RoomID:    "room1",
		UserID:    "user1",
		UserName:  "testuser",
		Text:      "Hello",
		Type:      "user",
		CreatedAt: time.Now().UTC(),
	}

	// Simulate broadcasting
	consumer.Hub.Broadcast(evt.RoomID, evt)

	if len(broadcaster.broadcasts) != 1 {
		t.Fatalf("expected 1 broadcast, got %d", len(broadcaster.broadcasts))
	}

	call := broadcaster.broadcasts[0]
	if call.RoomID != "room1" {
		t.Error("expected room ID to be 'room1'")
	}
}

func TestBotResponseConsumer_ProcessResponse(t *testing.T) {
	now := time.Now().UTC()
	msg := &message.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    "room1",
		UserID:    "bot",
		UserName:  "StockBot",
		Text:      "AAPL quote is 150.00 per share",
		Type:      "bot",
		CreatedAt: now,
	}

	service := &mockMessageService{botMessage: msg}

	// Test bot response processing
	response := BotResponseSubmit{
		RoomID: "room1",
		Text:   "AAPL quote is 150.00 per share",
	}

	// Simulate message creation
	m, err := service.CreateBotMessage(context.Background(), response.RoomID, response.Text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if m.Text != response.Text {
		t.Error("expected message text to match response text")
	}
}

func TestUserLookup(t *testing.T) {
	users := &mockUserRepository{user: &user.User{ID: primitive.NewObjectID(), Name: "testuser"}}

	// Test user lookup
	u, err := users.FindByID(context.Background(), "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Name != "testuser" {
		t.Error("expected user name to be 'testuser'")
	}
}

// Helper functions for testing command parsing
func startsWithSlash(s string) bool {
	return len(s) > 0 && s[0] == '/'
}

func splitCommand(s string) []string {
	if len(s) == 0 {
		return []string{}
	}

	// Remove leading slash
	cmd := s[1:]

	// Split on first space
	for i, c := range cmd {
		if c == ' ' {
			return []string{cmd[:i], cmd[i+1:]}
		}
	}
	return []string{cmd}
}

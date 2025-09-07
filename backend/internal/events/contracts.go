package events

import "time"

// Routing keys
const (
	RKMessageSubmit  = "message.submit"
	RKMessageCreated = "message.created"
)

type SubmitMessage struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
	Text   string `json:"text"`
}

type MessageCreated struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

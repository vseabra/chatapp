package events

import "time"

// Routing keys
const (
	RKMessageSubmit  = "message.submit"
	RKMessageCreated = "message.created"
	RKBotRequested   = "bot.requested"
	RKBotResponse    = "bot.response.submit"
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
	UserName  string    `json:"userName"`
	Text      string    `json:"text"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

type BotRequested struct {
	Command       string    `json:"command"`
	Args          string    `json:"args"`
	RoomID        string    `json:"roomId"`
	RequestUserID string    `json:"requestUserId"`
	MessageID     string    `json:"messageId"`
	RequestedAt   time.Time `json:"requestedAt"`
}

type BotResponseSubmit struct {
	RoomID string `json:"roomId"`
	Text   string `json:"text"`
}

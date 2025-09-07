package message

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents a chat message persisted in MongoDB.
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    string             `bson:"roomId" json:"roomId"`
	UserID    string             `bson:"userId" json:"userId"`
	UserName  string             `bson:"userName" json:"userName"`
	Text      string             `bson:"text" json:"text"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// ListResponse wraps paginated messages.
type ListResponse struct {
	Items      []Message `json:"items"`
	NextCursor string    `json:"nextCursor,omitempty"`
}

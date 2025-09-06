package chatroom

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatRoom represents a conversation room.
type ChatRoom struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Title     string             `bson:"title"`
	OwnerID   string             `bson:"ownerId"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"`
}

// CreateChatRoomRequest is the payload to create a chatroom.
type CreateChatRoomRequest struct {
	Title string `json:"title" binding:"required,min=1,max=120"`
}

// ChatRoomResponse is a public representation of a chatroom.
type ChatRoomResponse struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	OwnerID string `json:"ownerId"`
}

// UpdateChatRoomRequest allows renaming the chatroom.
type UpdateChatRoomRequest struct {
	Title string `json:"title" binding:"required,min=1,max=120"`
}

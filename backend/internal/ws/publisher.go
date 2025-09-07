package ws

import (
	"chatapp/internal/events"
	"context"
	"encoding/json"
)

type Publisher struct {
	AMQP *events.AMQP
}

func (p *Publisher) SubmitMessage(ctx context.Context, roomID string, userID string, text string) error {
	payload := events.SubmitMessage{RoomID: roomID, UserID: userID, Text: text}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.AMQP.PublishJSON(ctx, events.RKMessageSubmit, b)
}

package events

import (
	"chatapp/internal/message"
	"context"
	"encoding/json"
)

// IngressConsumer handles SubmitMessage, persists, emits MessageCreated.
type IngressConsumer struct {
	AMQP    *AMQP
	Service message.Service
}

func (c *IngressConsumer) Start(ctx context.Context) error {
	ch, err := c.AMQP.conn.Channel()
	if err != nil {
		return err
	}
	if _, err := ch.QueueDeclare("chat.ingress", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind("chat.ingress", RKMessageSubmit, c.AMQP.exchange, false, nil); err != nil {
		return err
	}
	msgs, err := ch.Consume("chat.ingress", "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var s SubmitMessage
			if err := json.Unmarshal(d.Body, &s); err != nil {
				continue
			}
			m, err := c.Service.Create(ctx, s.UserID, s.RoomID, s.Text)
			if err != nil {
				continue
			}
			evt := MessageCreated{ID: m.ID.Hex(), RoomID: m.RoomID, UserID: m.UserID, Text: m.Text, Type: m.Type, CreatedAt: m.CreatedAt}
			b, _ := json.Marshal(evt)
			_ = c.AMQP.PublishJSON(ctx, RKMessageCreated, b)
		}
	}()
	return nil
}

// BroadcastConsumer pushes MessageCreated events to the WebSocket hub.
type Broadcaster interface {
	Broadcast(roomID string, payload any)
}

type BroadcastConsumer struct {
	AMQP *AMQP
	Hub  Broadcaster
}

func (c *BroadcastConsumer) Start(ctx context.Context) error {
	ch, err := c.AMQP.conn.Channel()
	if err != nil {
		return err
	}
	if _, err := ch.QueueDeclare("chat.broadcast", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind("chat.broadcast", RKMessageCreated, c.AMQP.exchange, false, nil); err != nil {
		return err
	}
	msgs, err := ch.Consume("chat.broadcast", "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var evt MessageCreated
			if err := json.Unmarshal(d.Body, &evt); err != nil {
				continue
			}
			c.Hub.Broadcast(evt.RoomID, evt)
		}
	}()
	return nil
}

package events

import (
	"chatapp/internal/message"
	"chatapp/internal/user"
	"context"
	"encoding/json"
	"strings"
	"time"
)

// IngressConsumer handles SubmitMessage, persists, emits MessageCreated, and triggers bot requests.
type IngressConsumer struct {
	AMQP    *AMQP
	Service message.Service
	Users   user.Repository
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

			trim := strings.TrimSpace(s.Text)
			if strings.HasPrefix(trim, "/") {
				parts := strings.SplitN(trim[1:], " ", 2)
				cmd := parts[0]
				args := ""
				if len(parts) > 1 {
					args = parts[1]
				}
				req := BotRequested{Command: cmd, Args: args, RoomID: s.RoomID, RequestUserID: s.UserID, MessageID: "", RequestedAt: time.Now().UTC()}
				bb, _ := json.Marshal(req)
				_ = c.AMQP.PublishJSON(ctx, RKBotRequested, bb)

				// do not persist bot invocations to the DB
				continue
			}

			var userName string
			if c.Users != nil {
				if u, err := c.Users.FindByID(ctx, s.UserID); err == nil {
					userName = u.Name
				}
			}
			m, err := c.Service.CreateWithName(ctx, s.UserID, userName, s.RoomID, s.Text)
			if err != nil {
				continue
			}
			evt := MessageCreated{ID: m.ID.Hex(), RoomID: m.RoomID, UserID: m.UserID, UserName: m.UserName, Text: m.Text, Type: m.Type, CreatedAt: m.CreatedAt}
			b, _ := json.Marshal(evt)
			_ = c.AMQP.PublishJSON(ctx, RKMessageCreated, b)
		}
	}()
	return nil
}

// Broadcaster broadcasts to clients subscribed via WebSocket.
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

// BotResponseConsumer persists bot-origin messages and emits message.created.
type BotResponseConsumer struct {
	AMQP    *AMQP
	Service message.Service
}

func (c *BotResponseConsumer) Start(ctx context.Context) error {
	ch, err := c.AMQP.conn.Channel()
	if err != nil {
		return err
	}
	if _, err := ch.QueueDeclare("chat.bot.response", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind("chat.bot.response", RKBotResponse, c.AMQP.exchange, false, nil); err != nil {
		return err
	}
	msgs, err := ch.Consume("chat.bot.response", "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for d := range msgs {
			var r BotResponseSubmit
			if err := json.Unmarshal(d.Body, &r); err != nil {
				continue
			}
			m, err := c.Service.CreateBotMessage(ctx, r.RoomID, r.Text)
			if err != nil {
				continue
			}
			evt := MessageCreated{ID: m.ID.Hex(), RoomID: m.RoomID, UserID: m.UserID, UserName: m.UserName, Text: m.Text, Type: m.Type, CreatedAt: m.CreatedAt}
			b, _ := json.Marshal(evt)
			_ = c.AMQP.PublishJSON(ctx, RKMessageCreated, b)
		}
	}()
	return nil
}
